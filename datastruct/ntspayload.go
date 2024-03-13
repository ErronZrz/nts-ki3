package datastruct

import (
	"active/utils"
	"bytes"
	"fmt"
)

const (
	keyLength = 32
)

type NTSPayload struct {
	Host       string
	Port       int
	CertDomain string
	Secure     bool
	Err        error
	Len        int
	RcvData    []byte
	C2SKey     []byte
	S2CKey     []byte
}

type NTSDetectPayload struct {
	Host       string
	Port       int
	CertDomain string
	Secure     bool
	Info       *DetectInfo
}

type DetectInfo struct {
	CookieLength  int
	AEADList      []bool
	ServerPortSet map[string]struct{}
}

func (p *NTSPayload) Print() {
	if p.Err != nil {
		fmt.Println(p.Err)
	} else {
		fmt.Printf(p.Lines())
	}
}

func (p *NTSPayload) Lines() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("Remote address:     %s:%d (%s)\n", p.Host, p.Port, utils.RegionOf(p.Host)))
	if p.Secure {
		buf.WriteString(fmt.Sprintf("Certificate domain: %s (valid)\n", p.CertDomain))
	} else {
		if p.CertDomain != "" {
			buf.WriteString(fmt.Sprintf("Certificate domain: %s (unverified)\n", p.CertDomain))
		} else {
			buf.WriteString(fmt.Sprintf("Certificate not found"))
		}
	}
	if len(p.C2SKey) == keyLength {
		buf.WriteString("C2S key:            0x")
		for _, b := range p.C2SKey {
			buf.WriteString(fmt.Sprintf("%02X", b))
		}
		buf.WriteByte('\n')
	}
	if len(p.S2CKey) == keyLength {
		buf.WriteString("S2C key:            0x")
		for _, b := range p.S2CKey {
			buf.WriteString(fmt.Sprintf("%02X", b))
		}
		buf.WriteByte('\n')
	}
	buf.WriteString(fmt.Sprintf("%d bytes received:\n", p.Len))
	rows := p.Len >> 4
	for i := 0; i < rows; i++ {
		for _, b := range p.RcvData[i<<4 : (i+1)<<4] {
			buf.WriteString(fmt.Sprintf("%02X ", b))
		}
		buf.WriteByte('\n')
	}
	if p.Len > rows<<4 {
		for _, b := range p.RcvData[rows<<4:] {
			buf.WriteString(fmt.Sprintf("%02X ", b))
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (p *NTSDetectPayload) Lines() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("Remote address:     %s:%d (%s)\n", p.Host, p.Port, utils.RegionOf(p.Host)))
	if p.Secure {
		buf.WriteString(fmt.Sprintf("Certificate domain: %s (valid)\n", p.CertDomain))
	} else {
		if p.CertDomain != "" {
			buf.WriteString(fmt.Sprintf("Certificate domain: %s (unverified)\n", p.CertDomain))
		} else {
			buf.WriteString(fmt.Sprintf("Certificate not found"))
		}
	}

	supportNum := 0
	list := p.Info.AEADList
	for _, e := range list {
		if e {
			supportNum++
		}
	}
	buf.WriteString("Supported AEAD algorithms:")
	if supportNum == 0 {
		buf.WriteString(" None\n")
	} else {
		buf.WriteByte('\n')
		for id := byte(1); id <= 0x21; id++ {
			if list[id] {
				buf.WriteString(fmt.Sprintf("    - %s (%02X)\n", GetAEADName(id), id))
			}
		}
	}

	set := p.Info.ServerPortSet
	buf.WriteString("NTPv4 server addresses:")
	if len(set) == 0 {
		buf.WriteString(" Default\n")
	} else {
		buf.WriteByte('\n')
		for address := range set {
			buf.WriteString(fmt.Sprintf("    - %s\n", address))
		}
	}

	buf.WriteByte('\n')

	return buf.String()
}
