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
	C2SKeyMap map[byte][]byte
	CookieMap map[byte][][]byte
	Server    string
	Port      string
	T1        map[byte]time.Time
	T2        map[byte]time.Time
	T3        map[byte]time.Time
	T4        map[byte]time.Time
	RealT1    map[byte]time.Time
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
