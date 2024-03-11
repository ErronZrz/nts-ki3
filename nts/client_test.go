package nts

import (
	"active/parser"
	"fmt"
	"testing"
)

func TestDialNTSKE(t *testing.T) {
	payload, err := DialNTSKE("194.58.207.74", "sth2.nts.netnod.se", 0x0F)
	if err != nil {
		t.Error(err)
		return
	}

	if payload.Len == 0 {
		fmt.Println("empty response")
		return
	}
	payload.Print()

	res, err := parser.ParseNTSResponse(payload.RcvData)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Print(res.Lines())
	}
}
