package output

import (
	"active/async"
	"active/parser"
	"active/udpdetect"
	"active/utils"
	"fmt"
	"testing"
	"time"
)

func TestWriteToFile(t *testing.T) {
	cidr := "203.107.6.0/22"
	dataCh := udpdetect.DialNetworkNTP(cidr)
	if dataCh == nil {
		t.Error("dataCh is nil")
		return
	}

	seqNum := 0
	now := time.Now()
	for p, ok := <-dataCh; ok; p, ok = <-dataCh {
		err := p.Err
		if err != nil {
			t.Error(err)
			continue
		}
		header, err := parser.ParseHeader(p.RcvData)
		if err != nil {
			t.Error(err)
		} else {
			seqNum++
			WriteToFile(p.Lines(), header.Lines(), "test_timesync_"+cidr, seqNum, p.RcvTime, now)
		}
	}
	fmt.Printf("%d hosts detected in %s\n", seqNum, utils.DurationToStr(now, time.Now()))
}

func TestAsyncWriteToFile(t *testing.T) {
	cidr := "203.107.6.0/22"
	dataCh := async.DialNetworkNTP(cidr)
	if dataCh == nil {
		t.Error("dataCh is nil")
		return
	}

	seqNum := 0
	now := time.Now()
	for p, ok := <-dataCh; ok; p, ok = <-dataCh {
		err := p.Err
		if err != nil {
			t.Error(err)
			continue
		}
		header, err := parser.ParseHeader(p.RcvData)
		if err != nil {
			t.Error(err)
		} else {
			seqNum++
			WriteToFile(p.Lines(), header.Lines(), "test_timesync_"+cidr, seqNum, p.RcvTime, now)
		}
	}
	fmt.Printf("%d hosts detected in %s\n", seqNum, utils.DurationToStr(now, time.Now()))
}
