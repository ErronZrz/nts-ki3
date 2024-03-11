package dns

import "testing"

const (
	src    = "D:/Desktop/Detect/domain/final/EXTRA2.txt"
	ipDst  = "D:/Desktop/Detect/domain/final/EXTRA2_ip.txt"
	staDst = "D:/Desktop/Detect/domain/final/TWO2_sta.csv"
)

func TestOutputDNS(t *testing.T) {
	err := OutputDNS(src, ipDst)
	if err != nil {
		t.Error(err)
	}
}

func TestDetectAfterDNS(t *testing.T) {
	err := DetectAfterDNS(src, ipDst)
	if err != nil {
		t.Error(err)
	}
}

func TestAsyncDetectAfterDNS(t *testing.T) {
	err := AsyncDetectAfterDNS(src, ipDst)
	if err != nil {
		t.Error(err)
	}
}

func TestTLSAfterDNS(t *testing.T) {
	err := TLSAfterDNS(src, ipDst)
	if err != nil {
		t.Error(err)
	}
}

func TestDetectAEADAfterDNS(t *testing.T) {
	err := DetectAEADAfterDNS(src, ipDst)
	if err != nil {
		t.Error(err)
	}
}

func TestDetectStatisticAfterNTS(t *testing.T) {
	err := DetectStatisticAfterNTS(src, ipDst, staDst)
	if err != nil {
		t.Error(err)
	}
}

func TestAsyncDetectStatisticAfterDNS(t *testing.T) {
	err := AsyncDetectStatisticAfterDNS(src, ipDst, staDst)
	if err != nil {
		t.Error(err)
	}
}
