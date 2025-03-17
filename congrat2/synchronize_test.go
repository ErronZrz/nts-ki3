package congrat2

import (
	"active/congrat1"
	"database/sql"
	"testing"
)

func TestSynchronizeOnce(t *testing.T) {
	congrat1.UseDBConnection(func(db *sql.DB) error {
		return SynchronizeOnce(db, 30, 5, 5)
	})
}
