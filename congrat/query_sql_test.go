package congrat

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestMaxID(t *testing.T) {
	db, err := sql.Open("mysql", "root:liuyilun134@tcp(127.0.0.1:3306)/nts?charset=utf8")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = db.Close()
	}()

	maxID, err := MaxID(db)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("MaxID: %d\n", maxID)
	}
}
