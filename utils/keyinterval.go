package utils

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func AnalyzeInterval(prefix, t1, t2 string) (map[string]int, error) {
	times, err := timesBetween(t1, t2)
	if err != nil {
		return nil, err
	}
	m := make(map[string][]string)
	n := len(times)
	var idx int

	for _, timeStr := range times {
		path := prefix + timeStr + "_nts_keyid.txt"
		err := recordIDs(m, path, idx, n)
		idx++
		if err != nil {
			return nil, err
		}
	}

	res := make(map[string]int)
	badIDs := make([][]string, 0)
	violateIPs := make([]string, 0)
	for ip, ids := range m {
		lower, upper := intervalRange(ids)
		if lower == -2 {
			violateIPs = append(violateIPs, ip)
		}
		if lower > upper && upper >= 0 {
			badIDs = append(badIDs, ids)
		}
		res[fmt.Sprintf("(%d, %d)", lower, upper)]++
	}

	fmt.Printf("Violate IPs: %v\n", violateIPs)

	for _, ids := range badIDs {
		fmt.Print("? ")
		printIDs(ids)
	}

	return res, nil
}

func SaveIntervalTo(prefix, dst, t1, t2 string) error {
	file, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	writer := bufio.NewWriter(file)

	times, err := timesBetween(t1, t2)
	if err != nil {
		return err
	}
	m := make(map[string][]string)
	n := len(times)
	var idx int

	for _, timeStr := range times {
		path := prefix + timeStr + "_nts_keyid.txt"
		err := recordIDs(m, path, idx, n)
		idx++
		if err != nil {
			return err
		}
	}

	wanted := func(a, b int) bool {
		return ((a >= 20 && a <= 28) && (b >= 20 && b <= 28)) || ((a >= 164 && a <= 172) && (b >= 164 && b <= 172))
	}
	for ip, ids := range m {
		lower, upper := intervalRange(ids)
		if wanted(lower, upper) {
			_, err = writer.WriteString(fmt.Sprintf("%s\t%d\t%d\n", ip, lower, upper))
			if err != nil {
				return err
			}
		}
	}
	return writer.Flush()
}

func CrossCompare(kePath, itvPath string) (map[string]int, error) {
	keFile, err := os.Open(kePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = keFile.Close() }()
	itvFile, err := os.Open(itvPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = itvFile.Close() }()

	keScanner := bufio.NewScanner(keFile)
	itvScanner := bufio.NewScanner(itvFile)
	cookieLenMap := make(map[string]string)
	cookieLenNumMap := make(map[string]int)
	res := make(map[string]int)

	for keScanner.Scan() {
		ke := strings.Split(keScanner.Text(), "\t")
		ip := ke[0]
		cookieLen := betweenParentheses(ke[4])
		if _, ok := cookieLenMap[ip]; !ok {
			cookieLenMap[ip] = cookieLen
			cookieLenNumMap[cookieLen]++
		}
	}

	fmt.Println(cookieLenNumMap)

	for itvScanner.Scan() {
		itv := strings.Split(itvScanner.Text(), "\t")
		ip := itv[0]
		cookieLen, ok := cookieLenMap[ip]
		if !ok {
			continue
		}
		itvStr := "24"
		if len(itv[1]) > 2 {
			itvStr = "168"
		}
		res[cookieLen+"-"+itvStr]++
	}

	return res, nil
}

func timesBetween(t1, t2 string) ([]string, error) {
	layout := "2006010215"
	start, err := time.Parse(layout, t1)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(layout, t2)
	if err != nil {
		return nil, err
	}

	var times []string
	for !start.After(end) {
		times = append(times, start.Format(layout))
		start = start.Add(time.Hour)
	}

	return times, nil
}

func recordIDs(m map[string][]string, path string, idx, n int) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		ip, id := parts[0], parts[1]
		ids, ok := m[ip]
		if !ok {
			m[ip] = make([]string, n)
			ids = m[ip]
		}
		ids[idx] = id
	}

	return nil
}

func intervalRange(ids []string) (int, int) {
	ids = simplifyIDs(ids)
	printIDs(ids)

	type indexedValue struct {
		index int
		value string
	}
	var firstIndexes []*indexedValue
	var lastIndexes []*indexedValue

	lowerBound := -1
	exists := make(map[string]bool)

	for i, id := range ids {
		if len(id) == 0 {
			continue
		}
		if len(firstIndexes) == 0 {
			exists[id] = true
			firstIndexes = append(firstIndexes, &indexedValue{i, id})
			lastIndexes = append(lastIndexes, &indexedValue{i, id})
			continue
		}
		lastOne := firstIndexes[len(firstIndexes)-1]
		if lastOne.value == id {
			lowerBound = max(lowerBound, i-lastOne.index)
			lastIndexes[len(lastIndexes)-1].index = i
		} else if exists[id] {
			return -2, -2
		} else {
			exists[id] = true
			firstIndexes = append(firstIndexes, &indexedValue{i, id})
			lastIndexes = append(lastIndexes, &indexedValue{i, id})
		}
	}

	n := len(firstIndexes)
	if n < 3 {
		return lowerBound, -1
	}

	upperBound := math.MaxInt
	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				upperBound = min(upperBound, firstIndexes[k].index-lastIndexes[i].index)
			}
		}
	}

	if upperBound == math.MaxInt {
		return lowerBound, -1
	}
	return lowerBound, upperBound
}

func simplifyIDs(ids []string) []string {
	simplified := []byte{'A'}
	m := make(map[string]string)

	n := len(ids)
	res := make([]string, n)
	for i, id := range ids {
		if len(id) == 0 {
			continue
		}
		s, ok := m[id]
		if !ok {
			s = string(simplified)
			m[id] = s
			simplified[0]++
			if simplified[0] > 'Z' {
				simplified[0] = 'A'
			}
		}
		res[i] = s
	}
	return res
}

func printIDs(ids []string) {
	fmt.Print("[")
	n := len(ids)
	for i, id := range ids {
		if len(id) == 0 {
			fmt.Print("-")
		} else {
			fmt.Print(id)
		}
		if i < n-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println("]")
}

func betweenParentheses(s string) string {
	// 查找第一个左括号的位置
	start := strings.Index(s, "(")
	if start == -1 {
		return "" // 没有找到左括号，返回空字符串
	}

	// 查找第一个右括号的位置
	end := strings.Index(s[start+1:], ")")
	if end == -1 {
		return "" // 没有找到右括号，返回空字符串
	}

	// 提取并返回括号内的字符串
	return s[start+1 : start+1+end]
}
