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

func insertServerInfo(db *sql.DB, ip string, info *datastruct.OffsetServerInfo) error {
	query := `INSERT INTO ke_servers (batch_id, ip_address, domain_name, cert_org, cert_issuer, 
        ntpv4_address, ntpv4_port, domain_matches_ip, cert_not_expired, cert_not_self_signed, 
        cert_not_before, cert_not_after, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := db.Exec(query, CurrentBatchID, ip, info.CommonName, info.Organization,
		info.Issuer, info.Server, info.Port, info.RightIP, !info.Expired,
		!info.SelfSigned, adjustTime(info.NotBefore), adjustTime(info.NotAfter))
	// fmt.Println(adjustTime(info.NotBefore), adjustTime(info.NotAfter))
	if err != nil {
		return fmt.Errorf("error inserting server info: %v", err)
	}
	return nil
}

func insertKeyTimestamps(db *sql.DB, ip string, info *datastruct.OffsetServerInfo) error {
	query := `INSERT INTO ke_key_timestamp (batch_id, ip_address, availability, score, aead_id, c2s_key, s2c_key, 
        cookies, packet_len, ttl, stratum, poll, ntp_precision, root_delay, root_dispersion, reference,
        t1, t1r, t2, t3, t4, created_at, updated_at)
        VALUES (?, ?, 0, 0, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	for aeadID, timestamps := range info.Timestamps {
		packetLen := info.PacketLen[aeadID]
		stratum := info.Strata[aeadID]
		poll := info.Polls[aeadID]
		precision := info.Precisions[aeadID]
		rootDelay := info.RootDelays[aeadID]
		rootDispersion := info.RootDispersions[aeadID]
		reference := info.References[aeadID]
		c2sKey := info.C2SKeyMap[aeadID]
		s2cKey := info.S2CKeyMap[aeadID]
		var cookies []byte
		cookieList := info.CookieMap[aeadID]
		if len(cookieList) > 0 {
			n := len(cookieList[0])
			cookies = make([]byte, n*len(cookieList))
			for i, cookie := range cookieList {
				copy(cookies[i*n:], cookie)
			}
		}
		var t1r []byte
		t1, t2, t3 := timestamps[:8], timestamps[8:16], timestamps[16:]
		t4 := utils.GetTimestamp(info.T4[aeadID])
		if aeadID == 0 {
			t1r = t1
		} else {
			t1r = utils.GetTimestamp(info.RealT1[aeadID])
		}

		_, err := db.Exec(query, CurrentBatchID, ip, aeadID, c2sKey, s2cKey, cookies, packetLen,
			stratum, poll, precision, rootDelay, rootDispersion, reference, t1, t1r, t2, t3, t4)
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
