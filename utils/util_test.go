package utils

import (
	"fmt"
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
	tm := time.Now().AddDate(10, 0, 0)
	timestamp := GetTimestamp(tm)
	fmt.Println(timestamp)
	fmt.Println(ParseTimestamp(timestamp))
}

func TestRootDelayToValue(t *testing.T) {
	fmt.Println(RootDelayToValue([]byte{0, 0, 6, 0}))
}
