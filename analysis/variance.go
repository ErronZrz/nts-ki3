package analysis

import (
	"encoding/csv"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

type varParams struct {
	valCol       int
	yText        string
	unit         string
	divisor      float64
	low, high    float64
	useGlobalAvg bool
	syncOnly     bool
}

type avgMedianVar struct {
	avg    []float64
	median []float64
	std    []float64
	labels []string
}

var (
	sharedVarParams varParams
	max             float64
	syncMax         float64
	globalAvg       float64
	containsZero    bool
	sharedData      avgMedianVar
)

func VarianceBarChart(srcPath, dstDir, prefix string, params varParams) error {
	sharedVarParams = params

	lists, err := getStratumLists(srcPath)
	if err != nil {
		return err
	}

	containsZero = len(lists[0]) > 0

	generateData(lists)

	err = generateVarianceBarChart(dstDir, prefix)
	if err != nil {
		return err
	}

	return nil
}

func getStratumLists(srcPath string) ([][]float64, error) {
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

	res := make([][]float64, stratumLimit+2)
	all := make([]float64, 0)
	syn := make([]float64, 0)
	reader := csv.NewReader(file)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read file %s error: %v", srcPath, err)
		}

		valStr := row[sharedVarParams.valCol]
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse float %s error: %v", valStr, err)
		}
		realVal := float64(val) / sharedVarParams.divisor
		if realVal < sharedVarParams.low || realVal > sharedVarParams.high {
			continue
		}

		all = append(all, realVal)

		stratum, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse int %s error: %v", row[3], err)
		}
		if stratum == stratumLimit {
			stratum = 0
		}

		if stratum > 0 {
			syn = append(syn, realVal)
		}

		res[stratum] = append(res[stratum], realVal)
	}

	res[stratumLimit] = all
	res[stratumLimit+1] = syn

	return res, nil
}

func generateData(lists [][]float64) {
	sharedData = avgMedianVar{
		avg:    make([]float64, 0),
		median: make([]float64, 0),
		std:    make([]float64, 0),
		labels: []string{allName, synName},
	}

	getSingle(lists[stratumLimit], false)
	getSingle(lists[stratumLimit+1], true)
	globalAvg = sharedData.avg[1]

	for i := 0; i < stratumLimit; i++ {
		if len(lists[i]) == 0 {
			continue
		}
		sharedData.labels = append(sharedData.labels, getStratumStr(i))
		getSingle(lists[i], i > 0)
	}
}

func generateVarianceBarChart(dstDir string, prefix string) error {
	labels := sharedData.labels
	avg := sharedData.avg
	median := sharedData.median
	std := sharedData.std
	if sharedVarParams.syncOnly {
		start := 2
		if containsZero {
			start = 3
		}
		labels = append([]string{synName}, labels[start:]...)
		avg = append([]float64{avg[1]}, avg[start:]...)
		median = append([]float64{median[1]}, median[start:]...)
		std = append([]float64{std[1]}, std[start:]...)
	}
	n := len(labels)

	p := plot.New()
	p.Title.Text = fmt.Sprintf("Statistics of %s for Each Stratum", sharedVarParams.yText)
	p.X.Label.Text = "Stratum"
	p.Y.Label.Text = fmt.Sprintf("Average, Median and Standard Deviation of %s (%s)",
		sharedVarParams.yText, sharedVarParams.unit)
	if !sharedVarParams.syncOnly {
		p.Y.Max = stretchMax(max, false)
	} else {
		p.Y.Max = stretchMax(syncMax, false)
	}
	p.Y.Tick.Marker = plot.ConstantTicks(getMarks(p.Y.Max))

	width := vg.Points(20)
	barsAvg, err := plotter.NewBarChart(plotter.Values(avg), width)
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}
	barsMedian, err := plotter.NewBarChart(plotter.Values(median), width)
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}
	barsStd, err := plotter.NewBarChart(plotter.Values(std), width)
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}

	barsAvg.LineStyle.Width = vg.Length(0)
	barsMedian.LineStyle.Width = vg.Length(0)
	barsStd.LineStyle.Width = vg.Length(0)

	barsAvg.Color = plotutil.Color(0)
	barsMedian.Color = plotutil.Color(1)
	barsStd.Color = plotutil.Color(2)

	//barsAvg.ShowValue = true
	//barsMedian.ShowValue = true
	//barsStd.ShowValue = true

	barsAvg.Offset = -width
	barsMedian.Offset = 0
	barsStd.Offset = width

	p.Add(barsAvg, barsMedian, barsStd)
	p.Legend.Add("Average", barsAvg)
	p.Legend.Add("Median", barsMedian)
	p.Legend.Add("Standard Deviation", barsStd)
	p.Legend.Top = true
	p.Legend.Left = true
	p.NominalX(labels...)

	chartWidth := (2 + (vg.Length(n) * 0.8)) * vg.Inch
	chartHeight := 5 * vg.Inch

	err = p.Save(chartWidth, chartHeight, fmt.Sprintf("%s/%svar.png", dstDir, prefix))
	if err != nil {
		return fmt.Errorf("save chart error: %v", err)
	}

	return nil
}

func getSingle(list []float64, isSync bool) {
	avg := stat.Mean(list, nil)
	sharedData.avg = append(sharedData.avg, avg)

	sort.Float64s(list)
	median := stat.Quantile(0.5, stat.Empirical, list, nil)
	sharedData.median = append(sharedData.median, median)

	var std float64
	if len(sharedData.avg) <= 2 || !sharedVarParams.useGlobalAvg {
		std = standardDevOf(list, avg)
	} else {
		std = standardDevOf(list, globalAvg)
	}
	sharedData.std = append(sharedData.std, std)

	max = maxOfFour(max, avg, median, std)
	if isSync {
		syncMax = maxOfFour(syncMax, avg, median, std)
	}
}

func standardDevOf(list []float64, avg float64) float64 {
	var res float64
	for _, v := range list {
		diff := v - avg
		res += diff * diff
	}
	if res == 0 {
		return 0
	}

	n := len(list) - 1
	if n == 0 {
		n = 1
	}
	variance := res / float64(n)
	return math.Sqrt(variance)
}

func maxOfFour(a, b, c, d float64) float64 {
	return math.Max(math.Max(a, b), math.Max(c, d))
}
