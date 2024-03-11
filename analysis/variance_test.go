package analysis

import "testing"

func TestVarianceBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/final/ALL_BCIC8_ONE0_TWO0.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "Z_"
	params := varParams{
		valCol:       7,
		yText:        "Offset",
		unit:         "μs",
		divisor:      float64(2<<32) / 1000000,
		low:          -500,
		high:         500,
		useGlobalAvg: true,
		syncOnly:     false,
	}
	err := VarianceBarChart(srcPath, dstDir, prefix, params)
	if err != nil {
		t.Error(err)
	}
}
