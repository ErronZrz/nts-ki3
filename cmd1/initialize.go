package cmd1

import (
	"active/congrat1"
	"active/congrat2"
	"database/sql"
	"github.com/spf13/cobra"
)

var (
	m, minCandidates, minSurvivors int
	initializeCmd                  = &cobra.Command{
		Use:   "initialize",
		Short: "Initialize",
		Long:  "Initialize",
		Run: func(cmd *cobra.Command, args []string) {
			congrat1.UseDBConnection(func(db *sql.DB) error {
				return congrat2.Initialize(db, m, minCandidates, minSurvivors)
			})
		},
	}
)

func init() {
	initializeCmd.Flags().IntVarP(&m, "m", "m", 40, "m")
	initializeCmd.Flags().IntVarP(&minCandidates, "minCandidates", "c", 5, "minCandidates")
	initializeCmd.Flags().IntVarP(&minSurvivors, "minSurvivors", "s", 5, "minSurvivors")
}
