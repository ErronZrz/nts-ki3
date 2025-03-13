package congrat2

import (
	"active/congrat1"
	"database/sql"
	"fmt"
)

func insertServerInfoSimple(db *sql.DB, ke *KeKeyTimestamp, reserved string) error {
	query := `INSERT INTO ke_servers (batch_id, ip_address, domain_name, cert_org, cert_issuer, 
        ntpv4_address, ntpv4_port, domain_matches_ip, cert_not_expired, cert_not_self_signed, 
        cert_not_before, cert_not_after, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, TRUE, TRUE, TRUE, NOW(), NOW(), NOW(), NOW())`

	_, err := db.Exec(query, congrat1.CurrentBatchID, ke.IPAddress, reserved, reserved, reserved, ke.NTPv4Address, ke.NTPv4Port)
	// fmt.Println(adjustTime(ke.NotBefore), adjustTime(ke.NotAfter))
	if err != nil {
		return fmt.Errorf("error inserting server info: %v", err)
	}
	return nil
}

func insertKeyTimestamps2(db *sql.DB, ke *KeKeyTimestamp, aeadID int) error {
	query := `INSERT INTO ke_key_timestamp (batch_id, ip_address, availability, score, aead_id, c2s_key, s2c_key, 
        cookies, packet_len, ttl, stratum, poll, ntp_precision, root_delay, root_dispersion, reference,
        t1, t1r, t2, t3, t4, created_at, updated_at)
        VALUES (?, ?, 0, 0, ?, ?, ?, ?, ?, 63, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	if aeadID == 0 {
		ke.T1R = ke.T1
	}
	_, err := db.Exec(query, congrat1.CurrentBatchID, ke.IPAddress, aeadID, ke.C2SKey, ke.S2CKey,
		ke.Cookies, ke.PacketLen, ke.Stratum, ke.Poll, ke.NTPPrecision, ke.RootDelay, ke.RootDispersion,
		ke.Reference, ke.T1, ke.T1R, ke.T2, ke.T3, ke.T4)
	if err != nil {
		return fmt.Errorf("error inserting key timestamp for aeadID %d: %v", aeadID, err)
	}
	return nil
}
