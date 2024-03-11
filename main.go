package main

import (
	"active/datastruct"
	"active/nts"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	// 定义命令行参数
	var date string
	var n int
	var maxGoroutines int
	flag.StringVar(&date, "date", "", "The date in format YYYY-MM-DD")
	flag.IntVar(&n, "n", 0, "The number to be used")
	flag.IntVar(&maxGoroutines, "maxgoroutines", 100, "The maximum number of goroutines to run concurrently") // 默认值为100
	flag.Parse()

	// 检查日期参数是否已提供
	if date == "" {
		log.Fatalf("Error: The --date parameter is required")
	}
	// 获取用户的主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error: Unable to get the home directory: %v", err)
	}

	// 构建输入文件和输出文件的路径
	inputFilePath := fmt.Sprintf("%s/.nts/%s_ntske_ip_%d.txt", homeDir, date, n)
	outputFilePath := fmt.Sprintf("%s/.nts/%s_ntske_%d.txt", homeDir, date, n)

	// 打开包含主机地址的文本文件
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// 创建一个切片来存储主机地址
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
	outputFile, errs := os.Create(outputFilePath)
	if errs != nil {
		log.Fatalf("Error creating file %s: %v", outputFilePath, errs)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)

	// 创建等待组和互斥锁
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// 限制同时运行的协程数量
	//const maxGoroutines = 1600
	guard := make(chan struct{}, maxGoroutines)

	// 遍历每个主机地址并并发进行探测
	for _, host := range hosts {
		wg.Add(1)

		// 在协程中处理每个主机
		go func(h string) {
			defer wg.Done()
			guard <- struct{}{}        // 获取令牌
			defer func() { <-guard }() // 释放令牌

			serverName := "" // 替换为远程主机的 ServerName（如果有）
			// 进行探测
			result, err := nts.DetectNTSServer(h, serverName)
			if err != nil {
				log.Printf("Error detecting NTS server at %s: %v", h, err)
				//mutex.Lock() // 加锁
				//resultStr := fmt.Sprintf("%s						not support NTS						\n", h)
				//_, errs := writer.WriteString(resultStr)
				//mutex.Unlock() // 解锁
				if errs != nil {
					log.Printf("Error writing to output file: %v", errs)
				}
			} else {
				// 检查 AEADList 是否至少支持一个算法
				supportedAlgorithmExists := false
				for _, supported := range result.Info.AEADList {
					if supported {
						supportedAlgorithmExists = true
						break
					}
				}
				fmt.Printf("NTS Server Detection Result for %s:\n", h)
				if supportedAlgorithmExists {
					mutex.Lock() // 加锁
					results := fmt.Sprintf("%s\t", h)
					_, errs := writer.WriteString(results)
					mutex.Unlock() // 提前解锁
					if errs != nil {
						log.Printf("Error writing to output file: %v", errs)
						return
					}
					var maxSupportedID int = -1

					for id, supported := range result.Info.AEADList {
						if supported && id > maxSupportedID {
							maxSupportedID = id
						}
					}

					for id, supported := range result.Info.AEADList {
						//name := datastruct.GetAEADName(byte(id)) + ":"
						//status := "x"
						//var lastIndex int // 用于跟踪最后一个支持的元素的索引
						if supported {
							mutex.Lock() // 加锁
							//status = "supported"
							//fmt.Printf("- (%02X) %-27s   %s\n", id, name, status)
							name := datastruct.GetAEADName(byte(id))
							//lastIndex = id // 更新最后一个支持的元素的索引
							if id < maxSupportedID {
								results := fmt.Sprintf("%s\t", name)
								_, errs := writer.WriteString(results)
								mutex.Unlock() // 提前解锁
								if errs != nil {
									log.Printf("Error writing to output file: %v", errs)
									return
								}
							} else {
								results := fmt.Sprintf("%s", name)
								_, errs := writer.WriteString(results)
								mutex.Unlock() // 提前解锁
								if errs != nil {
									log.Printf("Error writing to output file: %v", errs)
									return
								}
							}
							//mutex.Unlock() // 提前解锁
						}
					}
					for serverPort := range result.Info.ServerPortSet {
						results := fmt.Sprintf(serverPort)
						mutex.Lock() // 加锁
						_, errs := writer.WriteString("\t" + results)
						mutex.Unlock() // 提前解锁
						if errs != nil {
							log.Printf("Error writing to output file: %v", errs)
							return
						}
					}
					resultStr := "\n"
					mutex.Lock() // 加锁
					_, errss := writer.WriteString(resultStr)
					mutex.Unlock() // 解锁
					if errss != nil {
						log.Printf("Error writing to output file: %v", errs)
					}
					//writer.Flush()
					//writer.Flush()
					fmt.Println()
				}
			}
		}(host)
		writer.Flush()
	}

	// 等待所有协程完成
	wg.Wait()

	// 将缓冲区内容写入文件
	writer.Flush()
	fmt.Printf("Results written to /.nts/%s_ntske_%d.txt\n", date, n)
}
