package offset

import (
	"testing"
)

func TestCalculateIPOffset(t *testing.T) {
	input := "D:\\Desktop\\TMP\\Ntages\\Data\\0606-500-1.txt"
	output := "D:\\Desktop\\TMP\\Ntages\\Data\\0606-500-1_offset-3.txt"
	err := CalculateOffsetsAsync(input, output, 1000)
	if err != nil {
		t.Error(err)
	}
}
