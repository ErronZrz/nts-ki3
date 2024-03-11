package analysis

import "testing"

func TestCountry4StratumBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/final/ALL_BCIC8_ONE0_TWO0.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "Z_"
	err := Country4StratumBarChart(srcPath, dstDir, prefix, true)
	if err != nil {
		t.Error(err)
	}
}
