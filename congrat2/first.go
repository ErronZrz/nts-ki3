package congrat2

import (
	"active/clock"
	"active/congrat1"
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

var GlobalSystemClock = new(clock.SystemClock)

// KeKeyTimestamp 结构体定义
type KeKeyTimestamp struct {
	ID             int
	BatchID        int
	IPAddress      string
	ASN            sql.NullInt32
	Availability   float64
	Score          float64
	AEADID         int
	C2SKey         []byte
	S2CKey         []byte
	Cookies        []byte
	PacketLen      int
	TTL            int
	Stratum        int
	Poll           int
	NTPPrecision   int
	RootDelay      []byte
	RootDispersion []byte
	Reference      []byte
	T1             []byte
	T1R            []byte
	T2             []byte
	T3             []byte
	T4             []byte
	CreatedAt      string
	UpdatedAt      string
	NTPv4Address   string
	NTPv4Port      int
}

func Initialize(db *sql.DB, m0, truechimerNum, survivorNum int) error {
	// 0. 更新批次
	maxBatchID, err := congrat1.MaxBatchID(db)
	if err != nil {
		return err
	}
	congrat1.CurrentBatchID = maxBatchID + 1
	// 1. 从数据库读取可用服务器
	serverList, err := getLatestRecordsByAEADID(db, 15)
	if err != nil {
		return err
	}
	fmt.Println(len(serverList))
	for _, server := range serverList {
		fmt.Println(server)
	}
	// 2. 以分数作为概率筛选 m0 台服务器
	selected := selectRecordsByScoreProbability(serverList, m0)
	fmt.Println(len(selected))
	for _, server := range selected {
		fmt.Println(server)
	}
	// 3. 进行同步
	for _, server := range selected {
		err = queryPort(db, server)
		if err != nil {
			return err
		}
		err = insertServerInfo2(db, server)
		if err != nil {
			return err
		}
		err = executeNTP(db, server, 15)
		if err != nil {
			fmt.Println(err)
			// 如果失败则标记失败
			server.NTPv4Address = FailFlag
			continue
		}
		err = insertKeyTimestamps2(db, server, 15)
		if err != nil {
			return err
		}
	}
	// 4. 生成对等体信息
	peers := getPeers(selected)
	// 5. 选出 truechimers、聚类、合并
	whatsoever(peers, truechimerNum, survivorNum)
	return nil
}

// 查询满足 aead_id=15 的最新记录
func getLatestRecordsByAEADID(db *sql.DB, aeadID int) ([]*KeKeyTimestamp, error) {
	query := `
		SELECT * FROM ke_key_timestamp WHERE id IN (
		    SELECT MAX(id) FROM ke_key_timestamp WHERE aead_id = ? GROUP BY ip_address
		)
	`

	rows, err := db.Query(query, aeadID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var results []*KeKeyTimestamp

	for rows.Next() {
		var record KeKeyTimestamp
		err := rows.Scan(
			&record.ID, &record.BatchID, &record.IPAddress, &record.ASN, &record.Availability, &record.Score, &record.AEADID,
			&record.C2SKey, &record.S2CKey, &record.Cookies, &record.PacketLen, &record.TTL, &record.Stratum, &record.Poll,
			&record.NTPPrecision, &record.RootDelay, &record.RootDispersion, &record.Reference,
			&record.T1, &record.T1R, &record.T2, &record.T3, &record.T4,
			&record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &record)
	}

	return results, nil
}

// 选取 m 条记录，概率正比于 Score
func selectRecordsByScoreProbability(records []*KeKeyTimestamp, m int) []*KeKeyTimestamp {
	if len(records) == 0 || m <= 0 {
		return nil
	}
	if m > len(records) {
		m = len(records)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 计算前缀和
	n := len(records)
	prefixSums := make([]float64, n)
	prefixSums[0] = records[0].Score
	for i := 1; i < n; i++ {
		prefixSums[i] = prefixSums[i-1] + records[i].Score
	}
	totalScore := prefixSums[n-1]

	// 进行加权随机选择
	selected := make([]*KeKeyTimestamp, 0, m)
	selectedSet := make(map[int]bool)

	for len(selected) < m {
		r := rng.Float64() * totalScore
		// 使用二分查找选择记录
		index := sort.Search(n, func(i int) bool {
			return prefixSums[i] >= r
		})

		if !selectedSet[records[index].ID] {
			selected = append(selected, records[index])
			selectedSet[records[index].ID] = true
		}
	}

	return selected
}
