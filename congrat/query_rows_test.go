package congrat

import (
	"database/sql"
	"testing"
)

func TestFetchAllTimestampRows(t *testing.T) {
	UseDBConnection(func(db *sql.DB) error {
		rows, err := FetchAllTimestampRows(db)
		if err != nil {
			return err
		}
		for _, row := range rows {
			t.Logf("%+v", row)
		}
		return nil
	})
}
