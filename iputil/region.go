package iputil

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"net"
	"strings"
)

const (
	nullIP      = "0.0.0.0"
	nullFlag    = "未同步"
	unknownFlag = "未知地区"
	privateFlag = "内网地址"
)

var (
	searcher *xdb.Searcher
)

func init() {
	xdbBuf, err := xdb.LoadContentFromFile("C:/Corner/TMP/毕设/NTP/Ntage3/nts-detect-txt/resources/ip2region.xdb")
	if err != nil {
		fmt.Printf("load xdb file error:%v\n", err)
	}
	searcher, err = xdb.NewWithBuffer(xdbBuf)
}

func GetChineseRegion(ipStr string, level int) string {
	if level < 1 || level > 3 {
		level = 3
	}
	if ipStr == nullIP {
		return nullFlag
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return unknownFlag
	}
	if ip.IsPrivate() {
		return privateFlag
	}
	region, err := searcher.SearchByStr(ipStr)
	if err != nil {
		fmt.Println(err)
		return unknownFlag
	}
	parts := strings.Split(region, "|")
	country := parts[0]
	if country == "0" {
		return unknownFlag
	}
	if country != "中国" {
		return "其他"
	}
	// 只有国家名或 level 为 1，则返回国家名
	if parts[2] == "0" || level == 1 {
		return country
	}
	// 直辖市或特别行政区，直接返回
	if strings.HasPrefix(parts[3], parts[2]) {
		return parts[2]
	}
	res := strings.ReplaceAll(parts[2], "省", "")
	// level 为 2，则返回省
	if parts[3] == "0" || level == 2 {
		return res
	}
	return res + strings.ReplaceAll(parts[3], "市", "")
}

func GetCountry(ipStr string) string {
	if ipStr == nullIP {
		return nullFlag
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return unknownFlag
	}
	if ip.IsPrivate() {
		return privateFlag
	}
	region, err := searcher.SearchByStr(ipStr)
	if err != nil {
		fmt.Println(err)
		return unknownFlag
	}
	parts := strings.Split(region, "|")
	country := parts[0]
	if country == "0" {
		return unknownFlag
	}
	return country
}
