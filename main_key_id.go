package main

import (
	"active/offset"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// 定义命令行参数
	var date string
	var maxCoroutines int
	flag.StringVar(&date, "date", "", "The date in format YYYY-MM-DD")
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
	inputFilePath := fmt.Sprintf("%s/.nts/%s_ntske_all.txt", homeDir, date)
	outputFilePath := fmt.Sprintf("%s/.nts/%s-%d_nts_keyid.txt", homeDir, date, time.Now().Hour())

	// 执行任务
	err = offset.GetNTSKeyID(inputFilePath, outputFilePath, maxCoroutines)
	if err != nil {
		log.Fatalf("Error: Unable to get the NTS Key ID: %v", err)
	}
}
