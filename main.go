package main

import (
	"active/datastruct"
	"active/nts"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	timeout int
)

func main() {
	// 定义命令行参数
	var date string
	var index int
	var maxGoroutines int
	flag.StringVar(&date, "date", "", "The date in format YYYY-MM-DD")
	flag.IntVar(&index, "index", 0, "The index of the file to process")
	flag.IntVar(&maxGoroutines, "maxgoroutines", 100, "The maximum number of goroutines to run at once")
	flag.IntVar(&timeout, "timeout", 20, "The timeout in seconds for each NTS server detection")
	flag.Parse()

	// 检查日期参数是否已提供
	if date == "" {
		log.Fatalf("Error: Parameter `date` is required")
	}

	// 获取用户的主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error: Unable to get the home directory: %v", err)
	}

	// 构建输入文件和输出文件的路径
	inputFilePath := fmt.Sprintf("%s/.nts/%s_ntske_ip_%d.txt", homeDir, date, index)
	outputFilePath := fmt.Sprintf("%s/.nts/%s_ntske_%d.txt", homeDir, date, index)

	// 打开包含主机地址的文本文件
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer func() { _ = file.Close() }()

	// 用于存储主机地址
	var hosts []string

	// 使用 bufio.Scanner 逐行读取主机地址
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := strings.TrimSpace(scanner.Text())
		if host != "" {
			hosts = append(hosts, host)
		}
	}

	// 检查扫描过程中是否发生错误
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning file: %v", err)
	}

	// 创建输出文件
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", outputFilePath, err)
	}
	defer func() { _ = outputFile.Close() }()

	writer := bufio.NewWriter(outputFile)

	// 创建 WaitGroup 以等待所有协程完成
	var wg sync.WaitGroup
	// 创建互斥锁以在写入文件时保护共享资源
	var mutex sync.Mutex
	// 限制同时运行的协程数量
	limitChan := make(chan struct{}, maxGoroutines)

	// 遍历主机地址并检测 NTS 服务器
	for _, host := range hosts {
		wg.Add(1)
		limitChan <- struct{}{}
		go func(host string) {
			defer wg.Done()
			detectAndWriteNTSServer(host, writer, &mutex, limitChan)
		}(host)
	}

	// 等待所有协程完成
	wg.Wait()

	// 将缓冲区内容写入文件
	_ = writer.Flush()
	fmt.Printf("Results written to /.nts/%s_ntske_%d.txt\n", date, index)
}

func detectAndWriteNTSServer(ip string, writer *bufio.Writer, mutex *sync.Mutex, limitChan chan struct{}) {
	// 释放信号量
	defer func() { <-limitChan }()

	// 进行探测
	result, err := nts.DetectNTSServer(ip, "", timeout)
	if err != nil {
		// log.Printf("Error detecting NTS server at %s: %v", ip, err)
		return
	}

	// 检查 AEADList 支持的算法
	supportedIds := make([]int, 0)
	info := result.Info
	for id, supported := range info.AEADList {
		if supported {
			supportedIds = append(supportedIds, id)
		}
	}
	// 若不支持任何 AEAD 算法，则打印消息并返回
	if len(supportedIds) == 0 {
		// log.Printf("No supported AEAD algorithms detected for %s", ip)
		return
	}

	// fmt.Printf("NTS Server Detection Result for %s:\n", ip)

	size := writer.Buffered()

	mutex.Lock()
	defer mutex.Unlock()

	// 1. 打印 IP 地址
	_, err = writer.WriteString(ip)
	newSize := writer.Buffered()
	// 如果 newSize == size，表明没有打印成功
	if newSize == size {
		log.Printf("Did not actually write string: %s", ip)
		return
	}
	size = newSize
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		return
	}

	// 2. 打印证书域名
	_, err = writer.WriteString("\t" + result.CertDomain)
	newSize = writer.Buffered()
	if newSize == size {
		log.Printf("Did not actually write string: %s", result.CertDomain)
		return
	}
	size = newSize
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		return
	}

	// 3. 打印指定的 NTPv4 服务器主机名和端口号
	var server, port string
	for serverPort := range info.ServerPortSet {
		parts := strings.Split(serverPort, ":")
		server, port = parts[0], parts[1]
	}
	if server == "" {
		server = "Default"
	}
	if port == "" {
		port = "123"
	}
	_, err = writer.WriteString("\t" + server + "\t" + port)
	newSize = writer.Buffered()
	if newSize == size {
		log.Printf("Did not actually write string: %s:%s", server, port)
		return
	}
	size = newSize
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		return
	}

	// 4. 打印支持的 AEAD 算法
	names := make([]string, len(supportedIds))
	for i, id := range supportedIds {
		names[i] = datastruct.GetAEADName(byte(id))
	}
	namesStr := strings.Join(names, ",")
	_, err = writer.WriteString("\t" + namesStr)
	newSize = writer.Buffered()
	if newSize == size {
		log.Printf("Did not actually write string: %s", namesStr)
		return
	}
	size = newSize
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		return
	}

	// 5. 打印 Cookie 长度
	lengthStr := strconv.Itoa(info.CookieLength)
	_, err = writer.WriteString("\t" + lengthStr)
	newSize = writer.Buffered()
	if newSize == size {
		log.Printf("Did not actually write string: %s", lengthStr)
		return
	}
	size = newSize
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		return
	}

	// 6. 打印是否过期以及是否自签名，0 为未过期，1 为已过期，2 为未生效，0 为非自签名，1 为自签名
	var expireFlag, selfSignedFlag int
	now := time.Now()
	if now.After(result.NotAfter) {
		expireFlag = 1
	} else if now.Before(result.NotBefore) {
		expireFlag = 2
	}
	if result.SelfSigned {
		selfSignedFlag = 1
	}
	_, err = writer.WriteString(fmt.Sprintf("\t%d\t%d", expireFlag, selfSignedFlag))
	newSize = writer.Buffered()
	if newSize == size {
		log.Printf("Did not actually write string: %d,%d", expireFlag, selfSignedFlag)
		return
	}
	size = newSize
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		return
	}

	// 7. 打印有效期
	layout := "2006-01-02 15:04:05"
	notBeforeStr := result.NotBefore.Format(layout)
	notAfterStr := result.NotAfter.Format(layout)
	_, err = writer.WriteString("\t" + notBeforeStr + "\t" + notAfterStr)
	newSize = writer.Buffered()
	if newSize == size {
		log.Printf("Did not actually write string: %s,%s", notBeforeStr, notAfterStr)
		return
	}
	size = newSize
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		return
	}

	// 8. 打印当前时间
	_, err = writer.WriteString("\t" + now.Format(layout))
	newSize = writer.Buffered()
	if newSize == size {
		log.Printf("Did not actually write string: %s", now.Format(layout))
		return
	}
	size = newSize
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		return
	}

	// 9. 打印换行符
	_, err = writer.WriteString("\n")
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
	}
}
