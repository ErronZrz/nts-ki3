package congrat1

import (
	"active/datastruct"
	"database/sql"
	"fmt"
	"testing"
)

func TestMainFunction(t *testing.T) {
	path := "C:\\Corner\\TMP\\BisheData\\0313-everNTS-585.txt"
	maxGoroutines := 10
	err := MainFunction(path, maxGoroutines)
	if err != nil {
		t.Error(err)
		return
	}

	UseDBConnection(func(db *sql.DB) error {
		maxID, err := MaxID(db)
		maxBatchID, err := MaxBatchID(db)
		CurrentBatchID = maxBatchID + 1
		fmt.Println("maxID:", maxID)
		fmt.Println("currentBatchID:", CurrentBatchID)
		for ip, info := range datastruct.OffsetInfoMap {
			err = InsertServerInfo(db, ip, info)
			if err != nil {
				return err
			}
			err = InsertKeyTimestamps(db, ip, info)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
