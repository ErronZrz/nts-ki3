package region

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
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
		if got := GetChineseRegion(test.input, 3); got != test.want {
			t.Errorf("GetChineseRegion(%s) = %s", test.input, got)
		}
	}
}

func TestRegionOfFile(t *testing.T) {
	path := "C:\\Corner\\TMP\\BisheData\\2024-08-17_ntps_all.txt"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		_ = file.Close()
	}()
	scanner := bufio.NewScanner(file)
	m := make(map[string]int)
	for scanner.Scan() {
		ip := strings.Split(scanner.Text(), "\t")[0]
		m[GetChineseRegion(ip, 2)]++
	}
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	fmt.Println(m)
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, m[k])
	}
}
