package dns

import (
	"testing"
)

func TestHandleText(t *testing.T) {
	path := "D:/Desktop/TMP/毕设/NTP/第六阶段/official-server/StratumTwo.txt"
	ipPath := "D:/Desktop/TMP/毕设/NTP/第六阶段/official-server/single-ip.txt"
	domainPath := "D:/Desktop/TMP/毕设/NTP/第六阶段/official-server/domain-two.txt"
	err := HandleText(path, ipPath, domainPath)
	if err != nil {
		t.Error(err)
	}
}
