package congrat

import (
	"database/sql"
	"testing"
)

func TestFetchAllTimestampRows(t *testing.T) {
	db, err := sql.Open("mysql", "root:liuyilun134@tcp(127.0.0.1:3306)/nts?charset=utf8&parseTime=true")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = db.Close()
	}()

	rows, err := FetchAllTimestampRows(db)
	if err != nil {
		t.Error(err)
	} else {
		for _, row := range rows {
			t.Logf("%+v", row)
		}
	}
}
