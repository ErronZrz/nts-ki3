package congrat

import (
	"database/sql"
	"testing"
)

func TestUpdateTTLWithFile(t *testing.T) {
	db, err := sql.Open("mysql", "root:liuyilun134@tcp(127.0.0.1:3306)/nts?charset=utf8")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = db.Close()
	}()

	path := "C:\\Corner\\TMP\\NTPData\\1126-4.pcapng"
	err = UpdateTTLWithFile(path, db)
	if err != nil {
		t.Error(err)
	}
}
