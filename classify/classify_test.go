package classify

import (
	"active/parser"
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestClassifyNTPRequest(t *testing.T) {
	filePath := "C:\\Corner\\TMP\\BisheData\\2024-01-27_ntps_passive_v4_5.pcap"
	packets, err := FetchNTPPackets(filePath, -1)
	if err != nil {
		t.Error(err)
		return
	}
	result := make(map[string]int)

	recordPath := "C:\\Corner\\TMP\\BisheData\\0405-5.txt"
	recordFile, err := os.Create(recordPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(recordFile)

	writer := bufio.NewWriter(recordFile)
	limit := 1000
	var otherCount int

	for _, packet := range packets {
		p, err := ParseNTPPacket(packet)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		str := ClassifyNTPRequest(p)
		if otherCount < limit && str == "other" {
			otherCount++
			for i := 0; i < parser.HeaderLength; i++ {
				_, _ = writer.WriteString(fmt.Sprintf("%02X", packet[i]))
				if i%16 == 15 {
					_, _ = writer.WriteString("\n")
				} else {
					_, _ = writer.WriteString(" ")
				}
			}
			_, _ = writer.WriteString("\n")
		}
		result[str]++
	}

	var count int
	for _, v := range result {
		count += v
	}
	fmt.Println(count)
	fmt.Println(result)
}
