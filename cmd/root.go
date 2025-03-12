package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ntpdtc",
		Short: "Ntpdtc is a tool used to detect NTP devices.",
		Long: "Ntpdtc is a tool used to detect NTP devices. The attributes that can be detected include " +
			"the OS version, running service and version, NTP reference clock information, etc.",
	}
)

func init() {
	rootCmd.AddCommand(timeSyncCmd)
	rootCmd.AddCommand(asyncCmd)
	rootCmd.AddCommand(ntsCmd)
	rootCmd.AddCommand(ntsAlgoCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
}

func handleError(cmd *cobra.Command, args []string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "execute %s args: %v error: %v\n", cmd.Name(), args, err)
	os.Exit(1)
}
