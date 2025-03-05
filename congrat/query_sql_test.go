package congrat

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestMaxID(t *testing.T) {
	UseDBConnection(func(db *sql.DB) error {
		maxID, err := MaxID(db)
		if err != nil {
			return err
		}
		fmt.Printf("MaxID: %d\n", maxID)
		return nil
	})
}
