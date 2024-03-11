package tcp

import (
	"active/utils"
	"fmt"
	"testing"
)

var (
	reqBytes = []byte{
		0x80, 0x01, 0x00, 0x02, 0x00, 0x00, 0x80, 0x04, 0x00, 0x02, 0x00, 0x0F, 0x80, 0x00, 0x00, 0x00,
	}
)

func TestIsTLSEnabled(t *testing.T) {
	var tests = []struct {
		host       string
		port       int
		serverName string
		want       bool
	}{
		{"bilibili.com", 443, "bilibili.com", true},
		{"bilibili.com", 443, "baidu.com", false},
		{"bilibili.com", 444, "bilibili.com", false},
		{"baidu.com", 443, "", true},
		{"192.168.179.129", 4460, "nothing.com", false},
		{"194.58.207.74", 4460, "sth2.nts.netnod.se", true},
		{"194.58.207.74", 4460, "", true},
	}

	for _, test := range tests {
		if got := IsTLSEnabled(test.host, test.port, test.serverName); got != test.want {
			t.Errorf("IsTLSEnabled(%s, %d) = %t", test.host, test.port, got)
		}
	}
}

func TestWriteReadTLS(t *testing.T) {
	res, err := WriteReadTLS("194.58.207.74", 4460, "sth2.nts.netnod.se", reqBytes)
	if err != nil {
		t.Error(err)
		return
	}
	n := len(res)
	if n == 0 {
		fmt.Println("empty response")
		return
	}
	fmt.Printf("Received %d bytes:\n", n)
	fmt.Print(utils.PrintBytes(res, 16))
}
