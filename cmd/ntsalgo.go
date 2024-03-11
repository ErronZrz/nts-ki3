package cmd

import "github.com/spf13/cobra"

var (
	ntsAlgoCmd = &cobra.Command{
		Use:   "ntsalgo <host> [domain]",
		Short: "Probe the algorithms that the server supports",
		Long: "Use the 'ntpdtc ntsalgo' command to establish the NTS-KE handshake multiple times to the" +
			" specified remote host to determine which AEAD algorithms are supported by the host.",
		Run: func(cmd *cobra.Command, args []string) {
			err := executeNTSAlgo(cmd, args)
			if err != nil {
				handleError(cmd, args, err)
			}
		},
	}
)
