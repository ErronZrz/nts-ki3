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
	path := "C:\\Corner\\TMP\\毕设\\NTP\\Ntage11\\0618-数据\\0618-all_offset-7"
	err := ExtractOffsetValues(path+".txt", path+"-256C.txt", 4, 5)
	err = ExtractOffsetValues(path+".txt", path+"-384C.txt", 6, 7)
	err = ExtractOffsetValues(path+".txt", path+"-512C.txt", 8, 9)
	err = ExtractOffsetValues(path+".txt", path+"-256S.txt", 3, 5)
	err = ExtractOffsetValues(path+".txt", path+"-384S.txt", 3, 7)
	err = ExtractOffsetValues(path+".txt", path+"-512S.txt", 3, 9)
	if err != nil {
		t.Error(err)
	}
}
