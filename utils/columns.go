package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ValueCount(path string, columns []int) ([]map[string]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	// 初始化用于统计的 map 切片，每个元素对应一个列的统计结果
	countMaps := make([]map[string]int, len(columns))
	for i := range countMaps {
		countMaps[i] = make(map[string]int)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 分割每行的数据
		fields := strings.Split(scanner.Text(), "\t")
		// 遍历所有需要统计的列
		for i, col := range columns {
			if col < len(fields) {
				// 更新统计信息
				val := fields[col]
				// 如果 val 是日期，则只保留年和月
				if strings.Contains(val, "-") && strings.Contains(val, ":") {
					val = val[:7]
				}
				countMaps[i][val]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return countMaps, nil
}
