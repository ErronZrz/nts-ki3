package cmd1

import (
	"active/congrat1"
	"active/congrat2"
	"database/sql"
	"github.com/spf13/cobra"
)

var (
	useKalman bool

	synchronizeCmd = &cobra.Command{
		Use:   "synchronize",
		Short: "Synchronize",
		Long:  "Synchronize",
		Run: func(cmd *cobra.Command, args []string) {
			congrat1.UseDBConnection(func(db *sql.DB) error {
				return congrat2.SynchronizeOnce(db, m, minCandidates, minSurvivors, useKalman)
			})
		},
	}
)

func init() {
	synchronizeCmd.Flags().IntVarP(&m, "m", "m", 30, "m")
	synchronizeCmd.Flags().IntVarP(&minCandidates, "minCandidates", "c", 5, "minCandidates")
	synchronizeCmd.Flags().IntVarP(&minSurvivors, "minSurvivors", "s", 5, "minSurvivors")
	synchronizeCmd.Flags().BoolVarP(&useKalman, "useKalman", "k", true, "useKalman")
}
