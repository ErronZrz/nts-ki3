package cmd1

import (
	"active/congrat1"
	"active/datastruct"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	maxGoroutines int
	buildPoolCmd  = &cobra.Command{
		Use:   "buildpool <path>",
		Short: "Build pool",
		Long:  "Build pool",
		Run: func(cmd *cobra.Command, args []string) {
			executeBuildPool(cmd, args)
		},
	}
)

func init() {
	buildPoolCmd.Flags().IntVarP(&maxGoroutines, "grnum", "g", 5,
		"Num of goroutines")
}

func executeBuildPool(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}
	err := congrat1.MainFunction(args[0], maxGoroutines)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	congrat1.UseDBConnection(func(db *sql.DB) error {
		maxID, err := congrat1.MaxID(db)
		maxBatchID, err := congrat1.MaxBatchID(db)
		congrat1.CurrentBatchID = maxBatchID + 1
		fmt.Println("maxID:", maxID)
		fmt.Println("currentBatchID:", congrat1.CurrentBatchID)
		for ip, info := range datastruct.OffsetInfoMap {
			err = congrat1.InsertServerInfo(db, ip, info)
			if err != nil {
				return err
			}
			err = congrat1.InsertKeyTimestamps(db, ip, info)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
