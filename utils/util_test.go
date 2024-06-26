package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func TestRegionOf(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"203.107.6.88", "广东深圳"},
		{"58.48.26.154", "湖北武汉"},
		{"218.27.132.24", "吉林吉林"},
		{"180.149.130.16", "北京"},
		{"101.227.131.220", "上海"},
		{"220.182.8.7", "西藏日喀则"},
		{"203.186.145.250", "香港"},
		{"114.44.227.87", "台湾"},
		{"1.2.3.4", "美国"},
		{"100.107.25.114", "未知地区"},
		{"192.168.179.128", "内网地址"},
	}
	for _, test := range tests {
		if got := RegionOf(test.input); got != test.want {
			t.Errorf("RegionOf(%s) = %s", test.input, got)
		}
	}
}

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

func TestRegionOfFile(t *testing.T) {
	path := "D:\\Desktop\\TMP\\Ntages\\Data\\2024-06-09_ntske_all.txt"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}
	scanner := bufio.NewScanner(file)
	m := make(map[string]int)
	for scanner.Scan() {
		ip := strings.Split(scanner.Text(), "\t")[0]
		m[RegionOf(ip)]++
	}
	for k, v := range m {
		fmt.Printf("%s: %d\n", k, v)
	}
}
