package congrat1

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// TimestampRow 表示数据库中的 ke_key_timestamp 表
type TimestampRow struct {
	ID             int
	IPAddress      string
	AEADID         int
	C2SKey         []byte
	S2CKey         []byte
	Cookies        []byte
	PacketLen      int
	TTL            int
	Stratum        int
	Poll           int
	Precision      int
	RootDelay      []byte
	RootDispersion []byte
	T1             []byte
	T2             []byte
	T3             []byte
	T4             []byte
	RealT1         []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// FetchAllTimestampRows 从数据库中查询所有记录
func FetchAllTimestampRows(db *sql.DB) ([]*TimestampRow, error) {
	rows, err := db.Query("SELECT * FROM ke_key_timestamp")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var timestamps []*TimestampRow
	for rows.Next() {
		var tr TimestampRow
		err := rows.Scan(&tr.ID, &tr.IPAddress, &tr.AEADID, &tr.C2SKey, &tr.S2CKey, &tr.Cookies,
			&tr.PacketLen, &tr.TTL, &tr.Stratum, &tr.Poll, &tr.Precision, &tr.RootDelay, &tr.RootDispersion,
			&tr.T1, &tr.RealT1, &tr.T2, &tr.T3, &tr.T4, &tr.CreatedAt, &tr.UpdatedAt)
		if err != nil {
			return nil, err
		}
		timestamps = append(timestamps, &tr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return timestamps, nil
}
