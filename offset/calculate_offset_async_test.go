package offset

import (
	"active/datastruct"
	"testing"
)

func TestCalculateIPOffset(t *testing.T) {
	input := "C:\\Corner\\TMP\\NTPData\\0530-all-1.txt"
	output := "C:\\Corner\\TMP\\NTPData\\0530-all-1_offset.txt"
	datastruct.OffsetInfoMap = make(map[string]*datastruct.OffsetServerInfo)
	err := CalculateOffsetsAsync(input, output)
	if err != nil {
		t.Error(err)
	}
}
