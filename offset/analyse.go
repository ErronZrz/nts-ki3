package offset

import (
	"bufio"
	"fmt"
	"os"
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
