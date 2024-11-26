package congrat

import (
	"active/datastruct"
	"active/utils"
	"database/sql"
	"fmt"
	"time"
)

var (
	mySQLStart = time.Date(1971, 1, 1, 0, 0, 1, 0, time.UTC)
	mySQLEnd   = time.Date(2037, 1, 19, 3, 14, 7, 0, time.UTC)
)

func insertServerInfo(db *sql.DB, ip string, serverInfo *datastruct.OffsetServerInfo) error {
	query := `INSERT INTO ke_servers (ip_address, domain_name, cert_org, cert_issuer, 
        ntpv4_address, ntpv4_port, domain_matches_ip, cert_not_expired, cert_not_self_signed, 
        cert_not_before, cert_not_after, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := db.Exec(query, ip, serverInfo.CommonName, serverInfo.Organization,
		serverInfo.Issuer, serverInfo.Server, serverInfo.Port, serverInfo.RightIP, !serverInfo.Expired,
		!serverInfo.SelfSigned, adjustTime(serverInfo.NotBefore), adjustTime(serverInfo.NotAfter))
	// fmt.Println(adjustTime(serverInfo.NotBefore), adjustTime(serverInfo.NotAfter))
	if err != nil {
		return fmt.Errorf("error inserting server info: %v", err)
	}
	return nil
}

func insertKeyTimestamps(db *sql.DB, ip string, serverInfo *datastruct.OffsetServerInfo) error {
	query := `INSERT INTO ke_key_timestamp (ip_address, aead_id, c2s_key, s2c_key, cookies, 
        packet_len, ttl, t1, t1r, t2, t3, t4, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	for aeadID, timestamps := range serverInfo.Timestamps {
		packetLen := serverInfo.PacketLen[aeadID]
		ttl := serverInfo.TTLs[aeadID]
		c2sKey := serverInfo.C2SKeyMap[aeadID]
		s2cKey := serverInfo.S2CKeyMap[aeadID]
		var cookies []byte
		cookieList := serverInfo.CookieMap[aeadID]
		if len(cookieList) > 0 {
			n := len(cookieList[0])
			cookies = make([]byte, n*len(cookieList))
			for i, cookie := range cookieList {
				copy(cookies[i*n:], cookie)
			}
		}
		var t1r []byte
		t1, t2, t3 := timestamps[:8], timestamps[8:16], timestamps[16:]
		t4 := utils.GetTimestamp(serverInfo.T4[aeadID])
		if aeadID == 0 {
			t1r = t1
		} else {
			t1r = utils.GetTimestamp(serverInfo.RealT1[aeadID])
		}

		_, err := db.Exec(query, ip, aeadID, c2sKey, s2cKey, cookies, packetLen, ttl, t1, t1r, t2, t3, t4)
		if err != nil {
			return fmt.Errorf("error inserting key timestamp for aeadID %d: %v", aeadID, err)
		}
	}
	return nil
}

func adjustTime(t time.Time) time.Time {
	if t.Before(mySQLStart) {
		return mySQLStart
	}
	if t.After(mySQLEnd) {
		return mySQLEnd
	}
	return t
}
