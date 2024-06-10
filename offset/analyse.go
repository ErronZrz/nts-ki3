package offset

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func AnalyseResult(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	m := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		ip, flag := parts[0], parts[2][3:4]
		already, ok := m[ip]
		if !ok {
			m[ip] = flag
		} else if already != flag {
			m[ip] = "?"
		}
	}

	for k, v := range m {
		fmt.Println(k + "\t" + v)
	}
	return nil
}

func AnalyseDomain(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	m := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		domain, flag := parts[1], parts[2][2:]
		m[domain] += flag
	}

	for k, v := range m {
		fmt.Println(k + "\t" + v)
	}
	return nil
}

func AnalyseOffset(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	res := make([][]float64, 7)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		for i, s := range parts[2:] {
			if len(s) == 1 {
				continue
			}
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			if f > -10000 && f < 10000 {
				res[i] = append(res[i], f)
			}
		}
	}

	for _, l := range res {
		var sum float64
		for _, v := range l {
			sum += v
		}
		fmt.Printf("len=%d, avg=%.6f\n", len(l), sum/float64(len(l)))
	}
	return nil
}
