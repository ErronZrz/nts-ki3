package validserver

import "testing"

func TestValidateKEServers(t *testing.T) {
	path := "D:\\Desktop\\TMP\\Ntages\\Ntage9\\0530-all.txt"
	path2 := "D:\\Desktop\\TMP\\Ntages\\Ntage9\\0530-all_offset.txt"
	err := ValidateKEServers(path, path2)
	if err != nil {
		t.Error(err)
	}
}
