package offset

import (
	"testing"
)

func TestCalculateIPOffset(t *testing.T) {
	input := "C:\\Corner\\TMP\\NTPData\\0530-all-1.txt"
	output := "C:\\Corner\\TMP\\NTPData\\0530-all-1_offset.txt"
	err := CalculateOffsetsAsync(input, output, 1000)
	if err != nil {
		t.Error(err)
	}
}
