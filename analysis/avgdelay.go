package analysis

import (
	"active/utils"
	"encoding/csv"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io"
	"os"
	"sort"
	"strconv"
)

type countryDelay struct {
	country  string
	eng      string
	avgDelay float64
}

func CountryAvgDelayBarChart(srcPath, dstDir, prefix string) error {
	mp, err := generateCountryDelayMap(srcPath)
	if err != nil {
		return err
	}

	list := mapToSortedList(mp)

	err = generateCountryAvgDelayBarChart(list, dstDir, prefix)
	if err != nil {
		return err
	}

	return nil
}

func generateCountryDelayMap(srcPath string) (map[string][]float64, error) {
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

	res := make(map[string][]float64)
	all := make([]float64, 0)
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
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv error: %v", err)
		}
		country := row[2]
		if _, ok := selected[country]; !ok {
			continue
		}

		delay, err := strconv.ParseInt(row[6], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse delay error: %v", err)
		}
		ms := float64(delay) / 1000
		all = append(all, ms)

		res[country] = append(res[country], ms)
	}

	res[global] = all

	return res, nil
}

func mapToSortedList(mp map[string][]float64) []countryDelay {
	list := make([]countryDelay, 0)

	for country, delays := range mp {
		cd := countryDelay{
			country:  country,
			eng:      "",
			avgDelay: stat.Mean(delays, nil),
		}
		list = append(list, cd)
	}

	for i, cd := range list {
		eng := utils.TranslateCountry([]string{cd.country})[0]
		list[i].eng = eng
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].avgDelay > list[j].avgDelay
	})

	return list
}

func generateCountryAvgDelayBarChart(list []countryDelay, dstDir, prefix string) error {
	n := len(list)
	values := make(plotter.Values, n)
	labels := make([]string, n)
	var max float64 = -1

	for i, cd := range list {
		values[i] = cd.avgDelay
		labels[i] = cd.eng
		if cd.avgDelay > max {
			max = cd.avgDelay
		}
	}

	p := plot.New()
	p.Title.Text = "Distribution of Average Delay for Each Country"
	p.Y.Label.Text = "Country"
	p.X.Label.Text = "Average Delay (ms)"
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
	p.NominalY(labels...)

	chartWidth := 6 * vg.Inch
	chartHeight := (1 + vg.Length(n)*0.3) * vg.Inch

	err = p.Save(chartWidth, chartHeight, fmt.Sprintf("%s/%scountry_avg_delay.png", dstDir, prefix))
	if err != nil {
		return fmt.Errorf("save chart error: %v", err)
	}

	return nil
}
