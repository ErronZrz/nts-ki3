package congrat

import (
	"database/sql"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
)

// 获取 ASN 信息
func getASN(ip string, db *geoip2.Reader) (uint, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ip)
	}

	record, err := db.ASN(parsedIP)
	if err != nil {
		return 0, err
	}

	return record.AutonomousSystemNumber, nil
}

// UpdateASN 更新数据库中的 ASN
func UpdateASN(db *sql.DB, geoDB *geoip2.Reader) error {
	// 查询所有 IP 地址
	rows, err := db.Query("SELECT id, ip_address FROM ke_key_timestamp WHERE asn IS NULL OR asn = 0")
	if err != nil {
		return err
	}
	defer func() {
		_ = rows.Close()
	}()

	m := make(map[string]uint)

	// 遍历查询结果
	var id int
	var ipAddress string
	var asn uint
	for rows.Next() {
		if err := rows.Scan(&id, &ipAddress); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		if asn, ok := m[ipAddress]; !ok {
			// 获取 ASN
			asn, err = getASN(ipAddress, geoDB)
			if err != nil {
				fmt.Println("Error getting ASN for IP:", ipAddress, err)
				continue
			}
			m[ipAddress] = asn
		}

		// 更新数据库
		_, err = db.Exec("UPDATE ke_key_timestamp SET asn = ? WHERE id = ?", asn, id)
		if err != nil {
			fmt.Println("Error updating ASN for ID:", id, err)
		} else {
			fmt.Printf("Updated IP %s with ASN %d\n", ipAddress, asn)
		}
	}

	return nil
}
