package classify

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	filePaths := []string{
		"C:\\Corner\\TMP\\BisheData\\1019-startover\\2024-09-14_mode6_0.pcap",
		"C:\\Corner\\TMP\\BisheData\\1019-startover\\2024-09-14_mode6_1.pcap",
		"C:\\Corner\\TMP\\BisheData\\1019-startover\\2024-09-14_mode6_2.pcap",
		"C:\\Corner\\TMP\\BisheData\\1019-startover\\2024-09-14_mode6_3.pcap",
	}
	dstFile, err := os.Create("C:\\Corner\\TMP\\BisheData\\1019-startover\\0914-mode6-system.txt")
	if err != nil {
		t.Error(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(dstFile)
	writer := bufio.NewWriter(dstFile)
	packets := make(map[string][][]byte)
	for _, filePath := range filePaths {
		err = FetchNTPPackets(filePath, packets)
		if err != nil {
			t.Error(err)
			return
		}
	}

	var total int

	for ip, packetList := range packets {
		total += len(packetList)
		for _, packet := range packetList {
			s := GetField(packet, "system")
			if s == "" {
				s = "no system"
			}
			_, _ = writer.WriteString(ip + "\t" + s + "\r\n")
		}
	}
	_ = writer.Flush()

	fmt.Println(total)
}

func TestClassifyNTPRequest(t *testing.T) {
	filePath := "C:\\Corner\\TMP\\BisheData\\2024-06-29_mode6_0.pcap"
	packets := make(map[string][][]byte)
	err := FetchNTPPackets(filePath, packets)
	if err != nil {
		t.Error(err)
		return
	}
	result := make(map[string]int)

	var tooShort int
	var total int

	for _, packetList := range packets {
		total += len(packetList)
		for _, packet := range packetList {
			s := GetField(packet, "version")
			isNtpd := strings.Contains(s, "ntpd")
			p, err := ParseNTPPacket(packet)
			if err != nil {
				tooShort++
				continue
			}
			str := ClassifyNTPRequest(p)
			if str == "ntpd" {
				if isNtpd {
					str = "ntpdT"
				} else {
					str = "ntpdF"
				}
			}
			result[str]++
		}
	}

	fmt.Println(total)
	fmt.Println(tooShort)

	var count int
	for _, v := range result {
		count += v
	}
	fmt.Println(count)
	fmt.Println(result)
}

func TestCross(t *testing.T) {
	filePath := "C:\\Corner\\TMP\\BisheData\\0629-1.txt"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	typeMap := make(map[string]int)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		ip := strings.Split(line, "\t")[0]
		if strings.Contains(line, "  20") {
			typeMap[ip] = 1
		} else if strings.Contains(line, "ntpd") {
			typeMap[ip] = 2
		} else {
			typeMap[ip] = 0
		}
	}

	cross("C:\\Corner\\TMP\\BisheData\\2024-01-27_ntps_passive_v4_0.pcap", typeMap)
	cross("C:\\Corner\\TMP\\BisheData\\2024-01-27_ntps_passive_v4_1.pcap", typeMap)
	cross("C:\\Corner\\TMP\\BisheData\\2024-01-27_ntps_passive_v4_2.pcap", typeMap)
	cross("C:\\Corner\\TMP\\BisheData\\2024-01-27_ntps_passive_v4_3.pcap", typeMap)
	cross("C:\\Corner\\TMP\\BisheData\\2024-01-27_ntps_passive_v4_4.pcap", typeMap)
	cross("C:\\Corner\\TMP\\BisheData\\2024-01-27_ntps_passive_v4_5.pcap", typeMap)
}

func cross(pcapPath string, typeMap map[string]int) {
	packets := make(map[string][][]byte)
	err := FetchNTPPackets(pcapPath, packets)
	if err != nil {
		return
	}

	var total int
	var tooShort int
	var not3 int
	res := make(map[string]int)
	for ip, packetList := range packets {
		total += len(packetList)
		flag, ok := typeMap[ip]
		for _, packet := range packetList {
			p, err := ParseNTPPacket(packet)
			if err != nil {
				tooShort++
				continue
			} else if p.Mode != 3 {
				not3++
				continue
			}
			s := ClassifyNTPRequest(p)
			if ok {
				s = fmt.Sprintf("%s %d", s, flag)
			}
			res[s]++
		}
	}
	fmt.Println(total)
	fmt.Println(not3)
	fmt.Println(tooShort)
	var count int
	for _, v := range res {
		count += v
	}
	fmt.Println(count)
	fmt.Println(res)
}
