package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func TestVariableData(t *testing.T) {
	want := time.Now().UTC().Format(timeFormat)
	data := VariableData()
	got := FormatTimestamp(data[40:])[:len(timeFormat)]
	if got != want {
		t.Errorf("VariableData timestamp: " + got)
	}
}

func TestSplitCIDR(t *testing.T) {
	got := SplitCIDR("192.168.254.147/22", 32)
	for _, s := range got {
		fmt.Println(s)
	}
}

func TestTranslateCountry(t *testing.T) {
	countries := []string{"中国", "德国", "美国", "加拿大"}
	result := TranslateCountry(countries)
	fmt.Println(result)
}

func TestParseNTPTimestamp(t *testing.T) {
	timestamp := make([]byte, 8)
	for i := range timestamp {
		timestamp[i] = 0xFF
	}
	parsed := ParseTimestamp(timestamp)
	fmt.Println(parsed)
}

func TestGetTimestamp(t *testing.T) {
	now := time.Now()
	t1 := now.Add(1 * time.Millisecond)
	t2 := t1.Add(50 * time.Microsecond)
	t3 := t2.Add(1 * time.Millisecond)
	t0 := GetTimestamp(now)
	str, err := bytesToHexList([][]byte{t0, t0, GetTimestamp(t1), GetTimestamp(t2), GetTimestamp(t3)})
	if err != nil {
		t.Errorf("Error converting bytes to hex: %v", err)
	} else {
		fmt.Println(str)
	}
}

func bytesToHexList(data [][]byte) (string, error) {
	var hexStrings []string

	for i, b := range data {
		if len(b) != 8 {
			return "", fmt.Errorf("element at index %d is not 8 bytes long (length=%d)", i, len(b))
		}
		hexStr := "0x" + strings.ToUpper(hex.EncodeToString(b))
		hexStrings = append(hexStrings, hexStr)
	}

	return strings.Join(hexStrings, ", "), nil
}

func TestRootDelayToValue(t *testing.T) {
	fmt.Println(RootDelayToValue([]byte{0, 0, 6, 0}))
}
