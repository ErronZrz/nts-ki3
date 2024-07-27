package offset

import "testing"

func TestGetNTSKeyID(t *testing.T) {
	inPath := "C:\\Users\\Jostle\\.nts\\2024-07-19_ntske_ip_2.txt"
	outPath := "C:\\Users\\Jostle\\.nts\\2024-07-19_ntske_keyid.txt"
	err := GetNTSKeyID(inPath, outPath, 100)
	if err != nil {
		t.Error(err)
	}
}
