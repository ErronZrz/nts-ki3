package iputil

import "testing"

func TestFindTimeDomains(t *testing.T) {
	path := "C:\\Corner\\TMP\\BisheData\\intersection_ntp_ip.txt"
	ptrPath := "C:\\Corner\\TMP\\BisheData\\intersection_ptr.txt"
	err := PTRRecords(path, ptrPath)
	if err != nil {
		t.Error(err)
	}
}
