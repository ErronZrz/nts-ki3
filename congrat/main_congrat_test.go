package congrat

import (
	"active/datastruct"
	"database/sql"
	"fmt"
	"testing"
)

func TestMainFunction(t *testing.T) {
	path := "C:\\Corner\\TMP\\BisheData\\1206-everNTS-100.txt"
	maxGoroutines := 20
	err := MainFunction(path, maxGoroutines)
	if err != nil {
		t.Error(err)
		return
	}
	db, err := sql.Open("mysql", "root:liuyilun134@tcp(127.0.0.1:3306)/nts?charset=utf8")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = db.Close()
	}()

	maxID, err := MaxID(db)
	fmt.Println("maxID:", maxID)
	for ip, info := range datastruct.OffsetInfoMap {
		err = insertServerInfo(db, ip, info)
		if err != nil {
			t.Error(err)
			return
		}
		err = insertKeyTimestamps(db, ip, info)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
