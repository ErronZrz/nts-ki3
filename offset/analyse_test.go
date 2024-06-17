package offset

import "testing"

func TestAnalyseResult(t *testing.T) {
	path := "D:\\Desktop\\TMP\\Ntages\\Ntage10\\offset_101.txt"
	err := AnalyseResult(path)
	if err != nil {
		t.Error(err)
	}
}

func TestAnalyseDomain(t *testing.T) {
	path := "D:\\Desktop\\TMP\\Ntages\\Data\\0606-500-1_offset-2.txt"
	err := AnalyseDomain(path)
	if err != nil {
		t.Error(err)
	}
}

func TestAnalyseOffset(t *testing.T) {
	path := "C:\\Corner\\TMP\\NTPData\\0606\\offset_100.txt"
	path = "D:\\Desktop\\TMP\\Ntages\\Ntage10\\offset_100.txt"
	err := AnalyseOffset(path)
	if err != nil {
		t.Error(err)
	}
}

func TestOffsetValues(t *testing.T) {
	path := "D:\\Desktop\\TMP\\Ntages\\Ntage10\\offset_100.txt"
	dstPath := "D:\\Desktop\\TMP\\Ntages\\Ntage10\\offset_100_512S.txt"
	err := ExtractOffsetValues(path, dstPath, 2, 8)
	if err != nil {
		t.Error(err)
	}
}
