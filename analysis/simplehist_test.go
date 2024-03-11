package analysis

import "testing"

func TestHistogramChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain28_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "domain28_9_"
	params := histParams{
		nameCol: 3,
		valCol:  8,
		subject: "Stratum ",
		xText:   "Processing Time (ms)",
		divisor: float64(2<<16) / 1000,
	}
	err := HistogramChart(srcPath, dstDir, prefix, params)
	if err != nil {
		t.Error(err)
	}
}
