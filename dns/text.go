package dns

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func HandleText(path, ipPath, domainPath string) error {
	srcFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()
	ipFile, err := os.OpenFile(ipPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = ipFile.Close() }()
	domainFile, err := os.OpenFile(domainPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = domainFile.Close() }()

	scanner := bufio.NewScanner(srcFile)
	ipWriter := bufio.NewWriter(ipFile)
	domainWriter := bufio.NewWriter(domainFile)
	checkCode := true
	codeMap := make(map[string]struct{})
	exist := make(map[string]struct{})

	for scanner.Scan() {
		line := scanner.Text()
		strs := strings.Split(line, "\t")
		code, ip, domain := strs[0], strs[1], strs[len(strs)-1]
		_, ok := codeMap[code]
		if checkCode && !ok {
			codeMap[code] = struct{}{}
			_, _ = domainWriter.WriteString(fmt.Sprintf("\n# %s\n", code))
			if code == "US" {
				checkCode = false
			}
		}
		if len(domain) > 1 {
			_, ok = exist[domain]
			if ok {
				continue
			}
			exist[domain] = struct{}{}
			_, _ = domainWriter.WriteString(domain)
			_ = domainWriter.WriteByte('\n')
		} else if len(ip) > 1 {
			_, ok = exist[ip]
			if ok {
				continue
			}
			exist[ip] = struct{}{}
			_, _ = ipWriter.WriteString(ip)
			_ = ipWriter.WriteByte('\n')
		}
	}

	_ = domainWriter.Flush()
	_ = ipWriter.Flush()

	return nil
}
