package congrat

import (
	"active/utils"
	"database/sql"
)

func UpdateTTLWithFile(path string, db *sql.DB) error {
	data, err := utils.ExtractTTLsAsMap(path)
	if err != nil {
		return err
	}
	return updateTTL(data, db)
}

func updateTTL(data map[string][][]int, db *sql.DB) error {
	// 首先查询所有 ip 及其 ID
	query := `
	SELECT MAX(id) AS id, ip_address AS ip
	FROM ke_key_timestamp
	GROUP BY ip_address
	`

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	idMap := make(map[string]int)
	for rows.Next() {
		id, ip := 0, ""
		if err := rows.Scan(&id, &ip); err != nil {
			return err
		}
		idMap[ip] = id
	}

	query = `
UPDATE ke_key_timestamp SET ttl = ?
WHERE id >= ? AND ip_address = ? AND packet_len = ?
`
	for ip, pairs := range data {
		id := idMap[ip] - 3
		for _, pair := range pairs {
			if _, err := db.Exec(query, pair[1], id, ip, pair[0]); err != nil {
				return err
			}
		}
	}

	return nil
}
