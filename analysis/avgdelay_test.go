package analysis

import "testing"

func TestCountryAvgDelayBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/final/ALL_BCIC8_ONE0_TWO0.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "Z_"
	err := CountryAvgDelayBarChart(srcPath, dstDir, prefix)
	if err != nil {
		t.Error(err)
	}
}
