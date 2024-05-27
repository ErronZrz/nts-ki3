package utils

import "testing"

func TestExtractNTPPackets(t *testing.T) {
	dataPath := "C:\\Corner\\TMP\\NTPData\\0527-2.pcapng"
	dstDir := "C:\\Corner\\TMP\\NTPData\\0527-packets"
	err := ExtractNTPPackets(dataPath, dstDir)
	if err != nil {
		t.Error(err)
	}
}
