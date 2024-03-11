package cmd

import "github.com/spf13/cobra"

var (
	ntsCmd = &cobra.Command{
		Use:   "nts <host> [domain]",
		Short: "Send NTS-KE request and parse response",
		Long: "Use the 'ntpdtc nts' command to establish a TLS connection with the specified " +
			"remote host, initiate an nts-KE request and listen for the response.",
		Run: func(cmd *cobra.Command, args []string) {
			err := executeNTS(cmd, args)
			if err != nil {
				handleError(cmd, args, err)
			}
		},
	}
)
