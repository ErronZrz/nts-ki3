package cmd

import (
	"github.com/spf13/cobra"
)

var (
	asyncCmd = &cobra.Command{
		Use:   "async <cidr>",
		Short: "Asynchronously sends and receives time synchronization packets",
		Long: "The 'async' command has the same effect as the 'timesync' command, but " +
			"sends and receives packets asynchronously.",
		Run: func(cmd *cobra.Command, args []string) {
			err := executeAsync(cmd, args)
			if err != nil {
				handleError(cmd, args, err)
			}
		},
	}
)

func init() {
	asyncCmd.Flags().IntVarP(&nPrintedHosts, "print", "p", 3,
		"The number of hosts you want to print out the results, no more than 16.")
}
