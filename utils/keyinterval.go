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
