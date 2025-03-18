package main

import (
	"active/datastruct"
	"active/offset"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main1() {
	// 定义命令行参数
	var date string
	var index int
	var maxCoroutines int
	flag.StringVar(&date, "date", "", "The date in format YYYY-MM-DD")
	flag.IntVar(&index, "index", 0, "The index of the file to process")
	flag.IntVar(&maxCoroutines, "maxCoroutines", 100, "Interval between tasks in milliseconds")
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

	// 执行任务
	err = offset.NTSKEDetectorWithFlags(inputFilePath, maxCoroutines)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 追加到文件
	outputFile, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", outputFilePath, err)
	}
	defer func() { _ = outputFile.Close() }()

	writer := bufio.NewWriter(outputFile)

	for ip, info := range datastruct.OffsetInfoMap {
		line := offset.GenerateLine2(ip, info)
		_, err = writer.WriteString(line)
		if err != nil {
			log.Fatalf("Error writing to file %s: %v", outputFilePath, err)
		}
	}

	err = writer.Flush()
	if err != nil {
		log.Fatalf("Error flushing file %s: %v", outputFilePath, err)
	}
}
