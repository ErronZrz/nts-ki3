package congrat

import (
	"database/sql"
	"testing"
)

func TestUpdateTTLWithFile(t *testing.T) {
	UseDBConnection(func(db *sql.DB) error {
		path := "C:\\Corner\\TMP\\NTPData\\1203-2.pcapng"
		err := UpdateTTLWithFile(path, db)
		return err
	})
}
