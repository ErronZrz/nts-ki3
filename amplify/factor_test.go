package amplify

import "testing"

func TestAmplificationFactors(t *testing.T) {
	path := "C:\\Corner\\TMP\\BisheData\\2024-03-02_monlist_0.pcap"
	outPath := "C:\\Corner\\TMP\\BisheData\\2024-03-02_monlist.txt"
	err := AmplificationFactors(path, outPath)
	if err != nil {
		t.Error(err)
	}
}
