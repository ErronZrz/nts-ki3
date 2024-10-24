package analysis

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	title1 = "Asynchronous Detection Results with Different Parameters"
	// title2  = "Synchronous Detection Results with Different Parameters"
	title2  = "Detection Results with Different Address Generation Methods"
	xLabel1 = "Sending Interval (ms)"
	// xLabel2 = "Listening Timeout (s)"
	xLabel2 = "Address Generation"
	yLabel  = "Device Count / Time Spent (s)"
)

var (
	legends1 = []string{"1 group", "4 groups", "16 groups"}
	legends2 = []string{"size = 64", "size = 128", "size = 256"}
	labels1  = []string{
		"15 / 60 / 240",
		"25 / 100 / 400",
		"40 / 160 / 640",
	}
	labels2 = []string{
		"Sequential",
		"Random",
	}
)

func DrawBarChart(data1, data2 [][]float64, dstPath1, dstPath2 string) error {
	err := drawChart1(data1, dstPath1)
	if err != nil {
		return err
	}

	err = drawChart2(data2, dstPath2)
	if err != nil {
		return err
	}

	return nil
}

func drawChart1(data [][]float64, path string) error {
	p := plot.New()
	p.Title.Text = title1
	p.X.Label.Text = xLabel1
	p.Y.Label.Text = yLabel
	p.Y.Max = stretchMax(data[5][2], false)
	p.Y.Tick.Marker = plot.ConstantTicks(getMarks(p.Y.Max))

	width := vg.Points(20)
	offset := -3 * width

	for i, v := range data {
		bars, err := plotter.NewBarChart(plotter.Values(v), width)
		if err != nil {
			return err
		}
		bars.LineStyle.Width = vg.Length(0)
		bars.Color = plotutil.Color(3 + i/2)
		//bars.ShowValue = true
		bars.Offset = offset
		offset += width

		p.Add(bars)
		if i&1 == 0 {
			p.Legend.Add(legends1[i/2], bars)
		}
	}

	p.Legend.Top = true
	p.Legend.Left = true
	p.NominalX(labels1...)

	chartWidth := 8 * vg.Inch
	chartHeight := 5 * vg.Inch

	return p.Save(chartWidth, chartHeight, path)
}

func drawChart2(data [][]float64, path string) error {
	p := plot.New()
	p.Title.Text = title2
	p.X.Label.Text = xLabel2
	p.Y.Label.Text = yLabel
	p.Y.Max = 5500
	p.Y.Tick.Marker = plot.ConstantTicks(getMarks(p.Y.Max))

	width := vg.Points(20)
	offset := -3 * width

	for i, v := range data {
		bars, err := plotter.NewBarChart(plotter.Values(v), width)
		if err != nil {
			return err
		}
		bars.LineStyle.Width = vg.Length(0)
		bars.Color = plotutil.Color(3 + i/2)
		//bars.ShowValue = true
		bars.Offset = offset
		offset += width

		p.Add(bars)
		if i&1 == 0 {
			p.Legend.Add(legends2[i/2], bars)
		}
	}

	p.Legend.Top = true
	p.Legend.Left = true
	p.NominalX(labels2...)

	chartWidth := 5 * vg.Inch
	chartHeight := 5 * vg.Inch

	return p.Save(chartWidth, chartHeight, path)
}

func DrawNTS(path string) error {
	p := plot.New()
	p.Title.Text = "Distribution of NTS Servers"
	p.X.Label.Text = "Country"
	p.Y.Label.Text = "Device Count"
	p.Y.Max = 16
	p.Y.Tick.Marker = plot.ConstantTicks(getMarks(p.Y.Max))

	width := vg.Points(20)
	barsTotal, err := plotter.NewBarChart(plotter.Values([]float64{15, 12, 7, 6, 4, 4, 3, 2, 1, 1, 1, 1, 1, 1}), width)
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}
	bars104, err := plotter.NewBarChart(plotter.Values([]float64{7, 0, 3, 2, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0}), width)
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}
	bars3, err := plotter.NewBarChart(plotter.Values([]float64{5, 0, 3, 2, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0}), width)
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}

	barsTotal.LineStyle.Width = vg.Length(0)
	bars104.LineStyle.Width = vg.Length(0)
	bars3.LineStyle.Width = vg.Length(0)

	barsTotal.Color = plotutil.Color(0)
	bars104.Color = plotutil.Color(1)
	bars3.Color = plotutil.Color(2)

	//barsTotal.ShowValue = true
	//bars104.ShowValue = true
	//bars3.ShowValue = true

	barsTotal.Offset = -width
	bars104.Offset = 0
	bars3.Offset = width

	p.Add(barsTotal, bars104, bars3)
	p.Legend.Add("Total Count", barsTotal)
	p.Legend.Add("Use 104-byte Cookie", bars104)
	p.Legend.Add("Support 3 AEAD Algorithms", bars3)
	p.Legend.Top = true
	p.NominalX("America", "Sweden", "Germany", "Brazil", "Netherlands", "Switzerland",
		"Slovenia", "Canada", "Britain", "Russia", "France", "Italy", "Denmark", "Croatia")

	chartWidth := 13.2 * vg.Inch
	chartHeight := 5 * vg.Inch

	err = p.Save(chartWidth, chartHeight, fmt.Sprintf("%s/nts.tif", path))
	if err != nil {
		return fmt.Errorf("save chart error: %v", err)
	}

	return nil
}
