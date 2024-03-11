package nts

import (
	"fmt"
	"testing"
)

func TestDetectNTSServer(t *testing.T) {
	host := "194.58.207.69"
	serverName := "sth1.nts.netnod.se"
	payload, err := DetectNTSServer(host, serverName)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Print(payload.Lines())
}
