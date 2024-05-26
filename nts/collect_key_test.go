package nts

import "testing"

func TestCollectKeys(t *testing.T) {
	path := "C:\\Users\\Jostle\\Desktop\\2024-05-05_ntske_all.txt"
	err := CollectKeys(path)
	if err != nil {
		t.Error(err)
	}
}
