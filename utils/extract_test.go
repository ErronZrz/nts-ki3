package utils

import (
	"fmt"
	"testing"
)

func TestExtractNTPPackets(t *testing.T) {
	dataPath := "C:\\Corner\\TMP\\NTPData\\0527-2.pcapng"
	dstDir := "C:\\Corner\\TMP\\NTPData\\0527-packets"
	err := ExtractNTPPackets(dataPath, dstDir)
	if err != nil {
		t.Error(err)
	}
}

func TestExtractTTLs(t *testing.T) {
	dataPath := "C:\\Corner\\TMP\\NTPData\\1126-1.pcapng"
	dstPath := "C:\\Corner\\TMP\\NTPData\\1126-ttl.txt"
	err := ExtractTTLs(dataPath, dstPath)
	if err != nil {
		t.Error(err)
	}
}

func TestExtractTTLsAsMap(t *testing.T) {
	dataPath := "C:\\Corner\\TMP\\NTPData\\1126-1.pcapng"
	m, err := ExtractTTLsAsMap(dataPath)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(m)
}
