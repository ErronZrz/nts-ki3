package offset

import "testing"

func TestCalculateOffsets(t *testing.T) {
	path := "C:\\Users\\Jostle\\Desktop\\0528-all-1.txt"
	err := CalculateOffsets(path)
	if err != nil {
		t.Error(err)
	}
}
