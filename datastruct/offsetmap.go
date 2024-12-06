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
	C2SKeyMap       map[byte][]byte
	S2CKeyMap       map[byte][]byte
	CookieMap       map[byte][][]byte
	RightIP         bool
	Expired         bool
	CommonName      string
	SelfSigned      bool
	NotBefore       time.Time
	NotAfter        time.Time
	Organization    string
	Issuer          string
	Server          string
	Port            string
	PacketLen       map[byte]int
	TTLs            map[byte]int
	Strata          map[byte]int
	Polls           map[byte]int
	Precisions      map[byte]int
	RootDelays      map[byte][]byte
	RootDispersions map[byte][]byte
	T1              map[byte]time.Time
	T2              map[byte]time.Time
	T3              map[byte]time.Time
	T4              map[byte]time.Time
	Timestamps      map[byte][]byte
	RealT1          map[byte]time.Time
}

func NewOffsetServerInfo(ip string) *OffsetServerInfo {
	return &OffsetServerInfo{
		C2SKeyMap:       make(map[byte][]byte),
		S2CKeyMap:       make(map[byte][]byte),
		CookieMap:       make(map[byte][][]byte),
		Server:          ip,
		Port:            "123",
		PacketLen:       make(map[byte]int),
		TTLs:            make(map[byte]int),
		Strata:          make(map[byte]int),
		Polls:           make(map[byte]int),
		Precisions:      make(map[byte]int),
		RootDelays:      make(map[byte][]byte),
		RootDispersions: make(map[byte][]byte),
		T1:              make(map[byte]time.Time),
		T2:              make(map[byte]time.Time),
		T3:              make(map[byte]time.Time),
		T4:              make(map[byte]time.Time),
		Timestamps:      make(map[byte][]byte),
		RealT1:          make(map[byte]time.Time),
	}
}

func (info *OffsetServerInfo) ClearTimeStamps() {
	info.Lock()
	info.T2[0x00] = time.Time{}
	info.T2[0x0F] = time.Time{}
	info.T2[0x10] = time.Time{}
	info.T2[0x11] = time.Time{}
	info.Timestamps[0x00] = []byte{}
	info.Timestamps[0x0F] = []byte{}
	info.Timestamps[0x10] = []byte{}
	info.Timestamps[0x11] = []byte{}
	info.Unlock()
}
