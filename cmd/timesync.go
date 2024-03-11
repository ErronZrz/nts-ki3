package cmd

import (
	"github.com/spf13/cobra"
)

var (
	nGoroutines   int
	nPrintedHosts int
	timeSyncCmd   = &cobra.Command{
		Use:   "timesync <cidr>",
		Short: "Send time synchronization requests and parse responses",
		Long: "Use the 'ntpdtc timesync' command to send a time synchronization request to the " +
			"specified CIDR address and listen for the response.",
		Run: func(cmd *cobra.Command, args []string) {
			err := executeTimeSync(cmd, args)
			if err != nil {
				handleError(cmd, args, err)
			}
		},
	}
)

func init() {
	timeSyncCmd.Flags().IntVarP(&nGoroutines, "grnum", "g", 0,
		"Num of goroutines. Setting it to 0 means using the value in the configuration file.")
	timeSyncCmd.Flags().IntVarP(&nPrintedHosts, "print", "p", 3,
		"The number of hosts you want to print out the results, no more than 16.")
}
