package congrat

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestGetServerExtensionFieldsCost(t *testing.T) {
	db, err := sql.Open("mysql", "root:liuyilun134@tcp(127.0.0.1:3306)/nts?charset=utf8")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = db.Close()
	}()

	ipTimestamps, err := fetchRecords(db)
	if err != nil {
		t.Error(err)
		return
	}

	for ip, mp := range ipTimestamps {
		if mp[0] == nil {
			continue
		}
		t1, t2, t3, t4 := mp[0].T1, mp[0].T2, mp[0].T3, mp[0].T4
		for aeadID, ts := range mp {
			if aeadID == 0 {
				continue
			}
			ntsT3, ntsT4 := ts.T3, ts.T4
			delay, offset, cost := GetServerExtensionFieldsCost(t1, t2, t3, t4, ntsT3, ntsT4)
			delayDuration := (time.Duration(delay) * time.Second) >> 32
			offsetDuration := (time.Duration(offset) * time.Second) >> 32
			costDuration := (time.Duration(cost) * time.Second) >> 32
			fmt.Printf("ip: %s, aeadID: %d, delay: %v, offset: %v, cost: %v\n", ip, aeadID, delayDuration, offsetDuration, costDuration)
		}
	}
}