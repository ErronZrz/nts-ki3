package validserver

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func ValidateKEServers(path, path2 string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	file2, err := os.Open(path2)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
		_ = file2.Close()
	}()

	scanner, scanner2 := bufio.NewScanner(file), bufio.NewScanner(file2)

	ntpServers := make(map[string]bool)
	for scanner2.Scan() {
		ntpServers[strings.Split(scanner2.Text(), "\t")[0]] = true
	}

	var yesYes, yesNo, noYes, noNo []string

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		ip, commonName, server := fields[0], fields[1], fields[2]
		check1 := checkDNS(commonName, ip)
		check2 := ntpServers[ip] || ntpServers[server]
		if check1 {
			if check2 {
				yesYes = append(yesYes, ip)
			} else {
				yesNo = append(yesNo, ip)
			}
		} else {
			if check2 {
				noYes = append(noYes, ip)
			} else {
				noNo = append(noNo, ip)
			}
		}
	}

	fmt.Println("YY:")
	for _, s := range yesYes {
		fmt.Println(s)
	}
	fmt.Println("YN:")
	for _, s := range yesNo {
		fmt.Println(s)
	}
	fmt.Println("NY:")
	for _, s := range noYes {
		fmt.Println(s)
	}
	fmt.Println("NN:")
	for _, s := range noNo {
		fmt.Println(s)
	}

	return nil
}

func checkDNS(domain, ip string) bool {
	if strings.Contains(domain, "*") {
		return true
	}

	ips, err := net.LookupIP(domain)
	if err != nil {
		return false
	}

	for _, resolvedIP := range ips {
		if resolvedIP.String() == ip {
			return true
		}
	}

	return false
}
