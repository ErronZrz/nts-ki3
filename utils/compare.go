package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func CompareIps(path1, path2 string) {
	ips1 := readIPsFromFile(path1)
	ips2 := readIPsFromFile(path2)

	fmt.Println("IPs in", path1, "but not in", path2)
	for ip, domain := range ips1 {
		if _, found := ips2[ip]; len(ip) > 0 && !found {
			fmt.Printf("%s (%s)\n", ip, domain)
		}
	}
	fmt.Println()

	fmt.Println("IPs in", path2, "but not in", path1)
	for ip, domain := range ips2 {
		if _, found := ips1[ip]; len(ip) > 0 && !found {
			fmt.Printf("%s (%s)\n", ip, domain)
		}
	}
}

func readIPsFromFile(path string) map[string]string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	ips := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) > 0 {
			ips[fields[0]] = fields[1]
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return ips
}
