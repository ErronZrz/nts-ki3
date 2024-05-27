package nts

import "testing"

func TestCollectKeys(t *testing.T) {
	path := "C:\\Users\\Jostle\\Desktop\\0527-all.txt"
	err := CollectKeys(path)
	if err != nil {
		t.Error(err)
	}
}

func TestMakeSecureNTPRequests(t *testing.T) {
	path := "C:\\Users\\Jostle\\Desktop\\0527-all_keys-1.txt"
	err := MakeSecureNTPRequests(path)
	if err != nil {
		t.Error(err)
	}
}
