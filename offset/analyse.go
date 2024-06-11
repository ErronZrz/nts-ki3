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

func ExtractOffsetValues(path, dstPath string, col1, col2 int) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
		_ = dstFile.Close()
	}()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(dstFile)
	var sum float64
	var count int

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts[col1]) == 1 || len(parts[col2]) == 1 {
			continue
		}
		f1, err := strconv.ParseFloat(parts[col1], 64)
		if err != nil {
			return err
		}
		f2, err := strconv.ParseFloat(parts[col2], 64)
		if err != nil {
			return err
		}
		sub := (f1 - f2) * 2
		if sub > -10000 && sub < 10000 {
			count++
			sum += sub
		}
		_, err = writer.WriteString(fmt.Sprintf("%.3f\n", sub))
		if err != nil {
			return err
		}
	}

	fmt.Printf("%.6f\n", sum/float64(count))

	return writer.Flush()
}
