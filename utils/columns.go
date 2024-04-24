package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
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

func PeriodCount(path string) (map[string]int, error) {
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

	countMap := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		countMap[monthStr(fields[8], fields[9])]++
	}

	return countMap, nil
}

func monthStr(d1, d2 string) string {
	// 定义时间格式
	layout := "2006-01-02 15:04:05"
	// 解析输入的时间字符串
	t1, _ := time.Parse(layout, d1)
	t2, _ := time.Parse(layout, d2)

	// 计算时间差异
	months := monthDiff(t1, t2)
	// 计算结果范围
	result := fmt.Sprintf("%d-%d", months, months+1)
	return result
}

func monthDiff(t1, t2 time.Time) int {
	years := t2.Year() - t1.Year()
	months := t2.Month() - t1.Month()
	days := t2.Day() - t1.Day()

	totalMonths := years*12 + int(months)
	if days < 0 {
		totalMonths--
	}
	return totalMonths
}
