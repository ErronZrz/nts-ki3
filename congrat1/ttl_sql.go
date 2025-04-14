package congrat1

import (
	"active/utils"
	"database/sql"
	"fmt"
	"math"
)

var initialTTLs = []int{64, 128, 255}

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

func UpdateAvailabilityAndScore(db *sql.DB) error {
	// 更新可用性
	query := `UPDATE ke_key_timestamp t1
		JOIN(
			SELECT ip_address, COUNT(*) as success_count FROM ke_key_timestamp WHERE aead_id = 15 GROUP BY ip_address
		) AS t1_count ON t1.ip_address = t1_count.ip_address
		JOIN(
			SELECT ip_address, COUNT(*) AS attempt_count FROM ke_servers GROUP BY ip_address
		) AS t2_count ON t1.ip_address = t2_count.ip_address
		SET t1.availability = t1_count.success_count / t2_count.attempt_count
		WHERE t1.availability = 0;
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error updating availability: %v", err)
	}

	// 更新分数
	query = `SELECT id, availability, stratum, root_delay, root_dispersion, ttl FROM ke_key_timestamp WHERE score = 0;`
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error selecting key timestamps: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var (
			id, stratum, ttl          int
			availability              float64
			rootDelay, rootDispersion []byte
		)
		err = rows.Scan(&id, &availability, &stratum, &rootDelay, &rootDispersion, &ttl)
		if err != nil {
			return fmt.Errorf("error scanning key timestamp: %v", err)
		}
		score := calculateScore(stratum, ttl, availability, rootDelay, rootDispersion)
		_, err = db.Exec("UPDATE ke_key_timestamp SET score = ? WHERE id = ?", score, id)
		if err != nil {
			return fmt.Errorf("error updating key timestamp: %v", err)
		}
	}

	return nil
}

func calculateScore(stratum, ttl int, availability float64, rootDelay, rootDispersion []byte) float64 {
	if stratum < 1 || stratum > 15 {
		return 50
	}
	base := 0.95
	availabilityFactor := math.Pow(base, 10*(1-availability))
	stratumFactor := math.Pow(base, 2*float64(stratum-1))
	rootDistance := utils.RootDelayToValue(rootDelay)/2 + utils.RootDelayToValue(rootDispersion)
	rootDistanceFactor := math.Pow(base, rootDistance*100)
	hops := 127
	for _, initialTTL := range initialTTLs {
		if ttl <= initialTTL {
			hops = initialTTL - ttl
		}
	}
	hopsFactor := math.Pow(base, float64(hops-1)/100)
	score := 100 * availabilityFactor * stratumFactor * rootDistanceFactor * hopsFactor
	fmt.Printf("ava = %.2f, str = %.2f, rd = %.2f, hop = %.2f, total = %.2f\n",
		availabilityFactor, stratumFactor, rootDistanceFactor, hopsFactor, score)
	return math.Max(50, score)
}
