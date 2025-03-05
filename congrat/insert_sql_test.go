package congrat

import (
	"active/datastruct"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

var info *datastruct.OffsetServerInfo

func init() {
	c2sKeys := [][]byte{
		{0x12, 0x34},
		{0x56, 0x78},
		{0x9a, 0xbc},
	}
	s2cKeys := [][]byte{
		{0xcb, 0xa9},
		{0x87, 0x65},
		{0x43, 0x21},
	}
	cookies := [][][]byte{
		{{0x12, 0x34}, {0x56, 0x78}},
		{{0x9a, 0xbc}, {0xcb, 0xa9}},
		{{0x87, 0x65}, {0x43, 0x21}},
	}
	t1s := []time.Time{
		time.Now().AddDate(0, 0, 1),
		time.Now().AddDate(0, 0, 2),
		time.Now().AddDate(0, 0, 3),
	}
	t1rs := []time.Time{
		time.Now().AddDate(0, 0, 4),
		time.Now().AddDate(0, 0, 5),
		time.Now().AddDate(0, 0, 6),
	}
	t2s := []time.Time{
		time.Now().AddDate(0, 0, 7),
		time.Now().AddDate(0, 0, 8),
		time.Now().AddDate(0, 0, 9),
	}
	t3s := []time.Time{
		time.Now().AddDate(0, 0, 10),
		time.Now().AddDate(0, 0, 11),
		time.Now().AddDate(0, 0, 12),
	}
	t4s := []time.Time{
		time.Now().AddDate(0, 0, 13),
		time.Now().AddDate(0, 0, 14),
		time.Now().AddDate(0, 0, 15),
	}
	info = &datastruct.OffsetServerInfo{
		Server:       "127.0.0.1",
		Port:         "4123",
		CommonName:   "test.ntp.com",
		RightIP:      true,
		Expired:      false,
		SelfSigned:   false,
		NotBefore:    time.Now().AddDate(-1, 0, 0),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		Organization: "Test Organization",
		Issuer:       "Test Issuer",
		C2SKeyMap:    make(map[byte][]byte),
		S2CKeyMap:    make(map[byte][]byte),
		CookieMap:    make(map[byte][][]byte),
		T1:           make(map[byte]time.Time),
		T2:           make(map[byte]time.Time),
		T3:           make(map[byte]time.Time),
		T4:           make(map[byte]time.Time),
		RealT1:       make(map[byte]time.Time),
	}
	aeads := []byte{0xab, 0xcd, 0xef}
	for i, id := range aeads {
		info.C2SKeyMap[id] = c2sKeys[i]
		info.S2CKeyMap[id] = s2cKeys[i]
		info.CookieMap[id] = cookies[i]
		info.T1[id] = t1s[i]
		info.RealT1[id] = t1rs[i]
		info.T2[id] = t2s[i]
		info.T3[id] = t3s[i]
		info.T4[id] = t4s[i]
	}
	// 改了时间戳赋值，所以这里不能这么写
	info.T1[0] = time.Now().AddDate(0, 1, 0)
	info.RealT1[0] = time.Now().AddDate(0, 2, 0)
	info.T2[0] = time.Now().AddDate(0, 3, 0)
	info.T3[0] = time.Now().AddDate(0, 4, 0)
	info.T4[0] = time.Now().AddDate(0, 5, 0)
}

func TestInsert(t *testing.T) {
	if info == nil {
		t.Error("info is nil")
		return
	}

	UseDBConnection(func(db *sql.DB) error {
		err := insertServerInfo(db, "127.0.0.1", info)
		if err != nil {
			return err
		}
		err = insertKeyTimestamps(db, "127.0.0.1", info)
		if err != nil {
			return err
		}
		return nil
	})
}

func TestAdjustTime(t *testing.T) {
	tm := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	fmt.Println(adjustTime(tm))
}
