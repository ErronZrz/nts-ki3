package analysis

import (
	"active/utils"
	"encoding/csv"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io"
	"os"
	"sort"
	"strconv"
)

type countryNum struct {
	country string
	eng     string
	num     int
}

const (
	specificSource = "特定时钟源"
)

func Country4StratumBarChart(srcPath, dstDir, prefix string, useRef bool) error {
	cnListList, err := generateCtr4SttSlice(srcPath, useRef)
	if err != nil {
		return err
	}

	c2eMap := make(map[string]string)
	for stratum, cnList := range cnListList {
		if len(cnList) == 0 {
			continue
		}

		sort.Slice(cnList, func(i, j int) bool {
			return cnList[i].num > cnList[j].num
		})

		for i, cn := range cnList {
			eng, ok := c2eMap[cn.country]
			if !ok {
				eng = utils.TranslateCountry([]string{cn.country})[0]
				c2eMap[cn.country] = eng
			}
			cnList[i].eng = eng
		}

		stratumStr := getStratumStr(stratum)

		err := generateCtr4SttBarChart(cnList, stratumStr, dstDir, prefix, useRef)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateCtr4SttSlice(srcPath string, useRef bool) ([][]countryNum, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("open file %s error: %v", srcPath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	res := make([][]countryNum, stratumLimit+2)
	all := make([]countryNum, 0)
	syn := make([]countryNum, 0)

	indexMap := make(map[string]int)
	reader := csv.NewReader(file)

	selected := map[string]struct{}{
		"中国":   {},
		"美国":   {},
		"日本":   {},
		"韩国":   {},
		"英国":   {},
		"新加坡":  {},
		"印度":   {},
		"南非":   {},
		"俄罗斯":  {},
		"德国":   {},
		"法国":   {},
		"意大利":  {},
		"西班牙":  {},
		"瑞士":   {},
		"巴西":   {},
		"澳大利亚": {},
		"加拿大":  {},
		"未知地区": {},
		"内网地址": {},
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv error: %v", err)
		}
		stratum, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse stratum error: %v", err)
		}
		if stratum >= stratumLimit {
			stratum = 0
		}
		now := res[stratum]

		country := row[2]
		if useRef {
			country = row[9]
		}
		if _, ok := selected[country]; !ok && country[0] >= 0x80 {
			continue
		}
		id := row[3] + country
		index, ok := indexMap[id]
		if ok {
			now[index].num++
		} else {
			indexMap[id] = len(now)
			now = append(now, countryNum{country: country, num: 1})
			res[stratum] = now
		}

		if useRef && country[0] < 0x80 {
			country = specificSource
		}
		id = "*" + country
		index, ok = indexMap[id]
		if ok {
			all[index].num++
		} else {
			indexMap[id] = len(all)
			all = append(all, countryNum{country: country, num: 1})
		}

		if stratum > 0 {
			id = "*" + id
			index, ok = indexMap[id]
			if ok {
				syn[index].num++
			} else {
				indexMap[id] = len(syn)
				syn = append(syn, countryNum{country: country, num: 1})
			}
		}
	}
	res[stratumLimit] = all
	res[stratumLimit+1] = syn

	return res, nil
}

func generateCtr4SttBarChart(cnList []countryNum, stratum, dstDir, prefix string, useRef bool) error {
	n := len(cnList)
	values := make(plotter.Values, n)
	var max float64 = -1
	for i := 0; i < n; i++ {
		// fmt.Println(cnList[i].eng, cnList[i].num)
		values[i] = float64(cnList[i].num)
		if values[i] > max {
			max = values[i]
		}
	}

	p := plot.New()
	if useRef {
		p.Title.Text = "Distribution of Ref Countries for " + stratum
	} else {
		p.Title.Text = "Distribution of Countries for " + stratum
	}
	p.Y.Label.Text = "Country"
	p.X.Label.Text = "Count"
	p.X.Max = stretchMax(max, true)
	p.X.Tick.Marker = plot.ConstantTicks(getMarks(p.X.Max))

	bars, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(0)
	bars.Horizontal = true
	bars.ShowValue = true

	p.Add(bars)
	xNames := make([]string, n)
	for i := 0; i < n; i++ {
		xNames[i] = cnList[i].eng
	}
	p.NominalY(xNames...)

	chartWidth := 4 * vg.Inch
	chartHeight := (1 + vg.Length(n)*0.3) * vg.Inch

	err = p.Save(chartWidth, chartHeight, fmt.Sprintf("%s/%s%s.png", dstDir, prefix, stratum))
	if err != nil {
		return fmt.Errorf("save chart error: %v", err)
	}

	return nil
}

func getStratumStr(stratum int) string {
	if stratum == 0 {
		return stratumNames[0]
	} else if stratum == stratumLimit {
		return allName
	} else if stratum == stratumLimit+1 {
		return synName
	}
	return "Stratum " + stratumNames[stratum]
}
