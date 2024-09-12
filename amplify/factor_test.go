package amplify

import "testing"

func TestAmplificationFactors(t *testing.T) {
	path := "C:\\Corner\\TMP\\BisheData\\08-28_reslist_0.pcap"
	outPath := "C:\\Corner\\TMP\\BisheData\\08-28_reslist.txt"
	err := AmplificationFactors(path, outPath)
	if err != nil {
		t.Error(err)
	}
}
