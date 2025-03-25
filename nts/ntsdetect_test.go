package nts

import (
	"fmt"
	"testing"
)

func TestDetectNTSServer(t *testing.T) {
	//host := "194.58.207.69"
	//serverName := "sth1.nts.netnod.se"
	host := "192.168.91.160"
	serverName := ""
	payload, err := DetectNTSServer(host, serverName, 20)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Print(payload.Lines())
}
