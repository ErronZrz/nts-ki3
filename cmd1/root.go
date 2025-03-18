package cmd1

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mynts",
		Short: "mynts is Erron's NTS.",
		Long:  "Too lazy to write description.",
	}
)

func init() {
	rootCmd.AddCommand(buildPoolCmd)
	rootCmd.AddCommand(updateTTLCmd)
	rootCmd.AddCommand(initializeCmd)
	rootCmd.AddCommand(synchronizeCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
}
