package offset

import "testing"

func TestCalculateOffsets(t *testing.T) {
	path := "D:\\Desktop\\TMP\\Ntages\\Ntage9\\0530-all.txt"
	err := CalculateOffsets(path)
	if err != nil {
		t.Error(err)
	}
}
