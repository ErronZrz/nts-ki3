package congrat

import (
	"active/datastruct"
	"active/utils"
	"database/sql"
	"fmt"
)

func insertServerInfo(db *sql.DB, serverInfo *datastruct.OffsetServerInfo) error {
	query := `INSERT INTO ke_servers (ip_address, domain_name, cert_org, cert_issuer, 
        ntpv4_address, ntpv4_port, domain_matches_ip, cert_not_expired, cert_not_self_signed, 
        cert_not_before, cert_not_after, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := db.Exec(query, serverInfo.Server, serverInfo.CommonName, serverInfo.Organization,
		serverInfo.Issuer, serverInfo.Server, serverInfo.Port, serverInfo.RightIP, !serverInfo.Expired,
		!serverInfo.SelfSigned, serverInfo.NotBefore, serverInfo.NotAfter)
	if err != nil {
		return fmt.Errorf("error inserting server info: %v", err)
	}
	return nil
}

func insertKeyTimestamps(db *sql.DB, serverInfo *datastruct.OffsetServerInfo) error {
	query := `INSERT INTO ke_key_timestamp (ip_address, aead_id, c2s_key, s2c_key, cookies, 
        t1, t1r, t2, t3, t4, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	for aeadID, timestamps := range serverInfo.Timestamps {
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

		_, err := db.Exec(query, serverInfo.Server, aeadID, c2sKey, s2cKey, cookies, t1, t1r, t2, t3, t4)
		if err != nil {
			return fmt.Errorf("error inserting key timestamp for aeadID %d: %v", aeadID, err)
		}
	}
	return nil
}
