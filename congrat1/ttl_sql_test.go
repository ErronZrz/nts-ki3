package congrat1

import (
	"database/sql"
	"testing"
)

func TestUpdateTTLWithFile(t *testing.T) {
	UseDBConnection(func(db *sql.DB) error {
		path := "C:\\Corner\\TMP\\NTPData\\0408-main.pcapng"
		err := UpdateTTLWithFile(path, db)
		if err != nil {
			return err
		}
		return UpdateAvailabilityAndScore(db)
	})
}

func TestUpdateJustScore(t *testing.T) {
	UseDBConnection(func(db *sql.DB) error {
		return UpdateAvailabilityAndScore(db)
	})
}
