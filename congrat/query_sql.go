package congrat

import "database/sql"

func MaxID(db *sql.DB) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT MAX(id) FROM ke_key_timestamp").Scan(&id)
	return id, err
}

func fetchRecords(db *sql.DB) (map[string]map[byte]*IPTimestamps, error) {
	query := `
	WITH RankedRecords AS (
	    SELECT id, ip_address, aead_id, t1, t1r, t2, t3, t4,
	           ROW_NUMBER() OVER (PARTITION BY ip_address, aead_id ORDER BY id DESC) AS rn
	    FROM ke_key_timestamp
	)
	SELECT id, ip_address, aead_id, t1, t1r, t2, t3, t4
	FROM RankedRecords
	WHERE rn = 1
	ORDER BY id;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	res := make(map[string]map[byte]*IPTimestamps)
	for rows.Next() {
		var kt IPTimestamps
		err = rows.Scan(&kt.ID, &kt.IPAddress, &kt.AeadID, &kt.T1, &kt.RealT1, &kt.T2, &kt.T3, &kt.T4)
		if err != nil {
			return nil, err
		}

		// 以IPAddress为key，将数据存入map中
		if res[kt.IPAddress] == nil {
			res[kt.IPAddress] = make(map[byte]*IPTimestamps)
		}
		res[kt.IPAddress][byte(kt.AeadID)] = &kt
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
