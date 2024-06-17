package offset

import (
	"active/nts"
	"testing"
)

func TestCalculateIPOffset(t *testing.T) {
	input := "C:\\Corner\\TMP\\NTPData\\0606\\0617-500-small.txt"
	output := "C:\\Corner\\TMP\\NTPData\\0606\\0617-500-small_offset.txt"
	nts.PlaceholderNum = 3
	err := CalculateOffsetsAsync(input, output, 1000)
	if err != nil {
		t.Error(err)
	}
}
