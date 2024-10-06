package classify

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFetchNTPPackets(t *testing.T) {
	pcapPath := "C:\\Corner\\TMP\\BisheData\\0914-mode6\\2024-09-14_mode6_0.pcap"
	keyMap := make(map[string]bool)
	data, err := FetchNTPPackets(pcapPath)
	if err != nil {
		t.Error(err)
		return
	}

	for _, d := range data {
		// 检查模式是否为 6，错误位是否为 0
		if len(d) <= 12 || d[0]&0x07 != 6 || d[1]&0x40 != 0 {
			continue
		}
		d = d[12:]
		kvs := bytes.Split(d, []byte{',', ' '})
		for _, kv := range kvs {
			i := bytes.Index(kv, []byte{'='})
			if i < 0 || i >= 20 {
				continue
			}
			keyMap[string(kv[:i])] = true
		}
	}
	for k := range keyMap {
		fmt.Println(k)
	}
}

func TestExtractMode6(t *testing.T) {
	pcapPath := "C:\\Corner\\TMP\\BisheData\\0914-mode6\\2024-09-14_mode6_0.pcap"
	outPath := "C:\\Corner\\TMP\\BisheData\\0914-mode6\\mode6_0.txt"
	data, err := FetchNTPPackets(pcapPath)
	if err != nil {
		t.Error(err)
		return
	}
	err = ExtractMode6(data, outPath)
	if err != nil {
		t.Error(err)
	}
}
