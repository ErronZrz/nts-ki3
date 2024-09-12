package region

import "testing"

func TestTopServers(t *testing.T) {
	path := "C:\\Corner\\TMP\\BisheData\\2024-08-17_ntps_all.txt"
	err := TopServers(path)
	if err != nil {
		t.Error(err)
	}
}
