package datastruct

import (
	"sync"
	"time"
)

var (
	OffsetInfoMap = make(map[string]*OffsetServerInfo)
	OffsetMapMu   sync.RWMutex
)

type OffsetServerInfo struct {
	sync.RWMutex
	C2SKeyMap    map[byte][]byte
	CookieMap    map[byte][][]byte
	RightIP      bool
	Expired      bool
	CommonName   string
	SelfSigned   bool
	NotBefore    time.Time
	NotAfter     time.Time
	Organization string
	Issuer       string
	Server       string
	Port         string
	T1           map[byte]time.Time
	T2           map[byte]time.Time
	T3           map[byte]time.Time
	T4           map[byte]time.Time
	RealT1       map[byte]time.Time
}

func NewOffsetServerInfo(ip string) *OffsetServerInfo {
	return &OffsetServerInfo{
		C2SKeyMap: make(map[byte][]byte),
		CookieMap: make(map[byte][][]byte),
		Server:    ip,
		Port:      "123",
		T1:        make(map[byte]time.Time),
		T2:        make(map[byte]time.Time),
		T3:        make(map[byte]time.Time),
		T4:        make(map[byte]time.Time),
		RealT1:    make(map[byte]time.Time),
	}
}

func (info *OffsetServerInfo) ClearTimeStamps() {
	info.Lock()
	info.T2[0x00] = time.Time{}
	info.T2[0x0F] = time.Time{}
	info.T2[0x10] = time.Time{}
	info.T2[0x11] = time.Time{}
	info.Unlock()
}
