package analysis

import (
	"encoding/csv"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io"
	"math"
	"os"
	"strconv"
)

var (
	minPrecision int64 = math.MaxInt64
	maxPrecision int64 = math.MinInt64
)

func StratumPrecisionBarChart(srcPath, dstDir, prefix string) error {
	spListList, err := generateStratumPrecisionSlice(srcPath)
	if err != nil {
		return err
	}

	for stratum, spList := range spListList {
		if len(spList) == 0 {
			continue
		}

		stratumStr := getStratumStr(stratum)

		err := generateStratumPrecisionBarChart(spList, stratumStr, dstDir, prefix)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateStratumPrecisionSlice(srcPath string) ([][]int64, error) {
	res := make([][]int64, stratumLimit+2)
	all := make([]int64, 0)
	syn := make([]int64, 0)

	file, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	reader := csv.NewReader(file)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read file %s error: %v", srcPath, err)
		}
		stratum, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse stratum error: %v", err)
		}
		if stratum >= stratumLimit {
			stratum = 0
		}

		precision, err := strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse precision error: %v", err)
		}
		all = append(all, precision)

		if stratum > 0 {
			syn = append(syn, precision)
		}

		if precision < minPrecision {
			minPrecision = precision
		}
		if precision > maxPrecision {
			maxPrecision = precision
		}

		now := res[stratum]
		now = append(now, precision)
		res[stratum] = now
	}

	res[stratumLimit] = all
	res[stratumLimit+1] = syn

	return res, nil
}

func generateStratumPrecisionBarChart(precisionList []int64, stratum, dstDir, prefix string) error {
	n := maxPrecision - minPrecision + 1
	values := make(plotter.Values, n)
	labels := make([]string, n)
	for _, precision := range precisionList {
		values[precision-minPrecision]++
	}
	var max float64 = 0
	for i, value := range values {
		if value > max {
			max = value
		}
		labels[i] = strconv.FormatInt(int64(i)+minPrecision, 10)
	}

	p := plot.New()
	p.Title.Text = fmt.Sprintf("Distribution of Poll for %s", stratum)
	p.X.Label.Text = "Poll"
	p.Y.Label.Text = "Count"
	p.Y.Max = stretchMax(max, false)
	p.Y.Tick.Marker = plot.ConstantTicks(getMarks(p.Y.Max))

	bars, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(0)
	bars.ShowValue = true

	p.Add(bars)
	p.NominalX(labels...)

	chartWidth := (1 + vg.Length(n)*0.3) * vg.Inch
	chartHeight := 4 * vg.Inch

	err = p.Save(chartWidth, chartHeight, fmt.Sprintf("%s/%s%s.png", dstDir, prefix, stratum))
	if err != nil {
		return fmt.Errorf("save bar chart error: %v", err)
	}

	return nil
}
