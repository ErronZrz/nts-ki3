package iputil

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type ProvinceServer struct {
	IP             string
	Stratum        int
	RootDispersion float64
	Precision      int
}

func TopServers(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	scanner := bufio.NewScanner(file)
	m := make(map[string][]*ProvinceServer)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		ip := parts[0]
		province := GetChineseRegion(ip, 2)
		if province == "其他" || province == "未知地区" {
			continue
		}
		stratum, _ := strconv.Atoi(parts[3])
		rootDispersion, _ := strconv.ParseFloat(parts[5], 64)
		precision, _ := strconv.Atoi(parts[2])
		ps := &ProvinceServer{
			IP:             ip,
			Stratum:        stratum,
			RootDispersion: rootDispersion,
			Precision:      precision,
		}
		m[province] = append(m[province], ps)
	}
	less := func(i, j *ProvinceServer) int {
		is, js := i.Stratum, j.Stratum
		if is == 0 {
			return 1
		} else if js == 0 {
			return -1
		}
		if is != js {
			return is - js
		}
		if i.RootDispersion < j.RootDispersion {
			return -1
		} else if i.RootDispersion > j.RootDispersion {
			return 1
		}
		return i.Precision - j.Precision
	}
	for _, v := range m {
		slices.SortFunc(v, less)
	}
	for k, v := range m {
		fmt.Println(k + "\n一层")
		var i int
		for i < 50 && v[i].Stratum == 1 {
			fmt.Print(v[i].IP + ", ")
			i++
		}
		for v[i].Stratum == 1 {
			i++
		}
		fmt.Println("\n二层")
		j := i + 50
		for i < j {
			fmt.Print(v[i].IP + ", ")
			i++
		}
		fmt.Println()
	}
	return nil
}
