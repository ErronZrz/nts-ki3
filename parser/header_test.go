package parser

import (
	"active/udpdetect"
	"fmt"
	"testing"
)

func TestParseHeader(t *testing.T) {
	dataCh := udpdetect.DialNetworkNTP("192.168.91.128/25")
	if dataCh == nil {
		fmt.Println("dataCh is nil")
		return
	}
	for p, ok := <-dataCh; ok; p, ok = <-dataCh {
		err := p.Err
		if err != nil {
			fmt.Println(err)
			continue
		}
		data := p.RcvData
		p.Print()
		header, err := ParseHeader(data)
		if err != nil {
			t.Error(err)
		} else {
			header.Print()
		}
	}
}
