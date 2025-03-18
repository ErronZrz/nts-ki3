package cmd1

import (
	"active/congrat1"
	"database/sql"
	"github.com/spf13/cobra"
)

var (
	updateTTLCmd = &cobra.Command{
		Use:   "updatettl <path>",
		Short: "Update TTL",
		Long:  `Update TTL`,
		Run: func(cmd *cobra.Command, args []string) {
			executeUpdateTTL(cmd, args)
		},
	}
)

func executeUpdateTTL(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}
	congrat1.UseDBConnection(func(db *sql.DB) error {
		err := congrat1.UpdateTTLWithFile(args[0], db)
		if err != nil {
			return err
		}
		return congrat1.UpdateAvailabilityAndScore(db)
	})
}
