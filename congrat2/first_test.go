package congrat2

import (
	"active/congrat1"
	"database/sql"
	"testing"
)

func TestInitialize(t *testing.T) {
	congrat1.UseDBConnection(func(db *sql.DB) error {
		return Initialize(db, 40, 5, 5)
	})
}
