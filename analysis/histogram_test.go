package analysis

import "testing"

func TestHistogramBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/final/ALL_BCIC8_ONE0_TWO0.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "Z_"
	partitions := []int{
		1, 3, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 60, 70, 80, 90, 100, 120, 150, 200, 300, 500, 1000,
	}

	params := histParams{
		nameCol:    3,
		valCol:     11,
		subject:    "Stratum",
		xText:      "Root Dispersion",
		unit:       "ms",
		divisor:    float64(2<<16) / 1000,
		partitions: partitions,
	}
	err := HistogramBarChart(srcPath, dstDir, prefix, params)
	if err != nil {
		t.Error(err)
	}
}
