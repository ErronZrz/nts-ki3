package cmd

import (
	"active/nts"
	"active/output"
	"active/parser"
	"active/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func executeNTS(cmd *cobra.Command, args []string) error {
	cmdName := cmd.Name()
	if args == nil || len(args) == 0 {
		return fmt.Errorf("command `%s` missing arguments", cmdName)
	}
	host := args[0]
	var serverName, serverNameLine string
	if len(args) > 2 {
		return fmt.Errorf("%d arguments in command `%s`, expecting 2", len(args), cmdName)
	}
	if len(args) == 2 {
		serverName = args[1]
		serverNameLine = fmt.Sprintf("    server name: %s\n", serverName)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Ready to run `%s`.\n    host: %s\n%s\n", cmdName, host, serverNameLine)

	startTime := time.Now()

	payload, err := nts.DialNTSKE(host, serverName, 0x0F)
	if err != nil {
		return err
	}
	if payload.Len == 0 {
		_, _ = fmt.Fprintln(os.Stdout, "Empty response.")
		return nil
	}

	res, err := parser.ParseNTSResponse(payload.RcvData)
	if err != nil {
		return err
	}

	raw, parsed := payload.Lines(), res.Lines()
	_, _ = fmt.Fprintf(os.Stdout, "%s[parsed]\n%s", raw, parsed)
	output.WriteNTSToFile(raw, parsed, host)

	_, _ = fmt.Fprintf(os.Stdout, "NTS-KE handshake and parsing were completed in %s\n",
		utils.DurationToStr(startTime, time.Now()))

	return nil
}

func executeNTSAlgo(cmd *cobra.Command, args []string) error {
	cmdName := cmd.Name()
	if args == nil || len(args) == 0 {
		return fmt.Errorf("command `%s` missing arguments", cmdName)
	}
	host := args[0]
	var serverName, serverNameLine string
	if len(args) > 2 {
		return fmt.Errorf("%d arguments in command `%s`, expecting 2", len(args), cmdName)
	}
	if len(args) == 2 {
		serverName = args[1]
		serverNameLine = fmt.Sprintf("    server name: %s\n", serverName)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Ready to run `%s`.\n    host: %s\n%s\n", cmdName, host, serverNameLine)

	startTime := time.Now()

	payload, err := nts.DetectNTSServer(host, serverName, 20)
	if err != nil {
		return err
	}

	content := payload.Lines()
	_, _ = fmt.Fprint(os.Stdout, content)
	output.WriteNTSDetectToFile(content, host)

	_, _ = fmt.Fprintf(os.Stdout, "33 algorithms detected in %s\n",
		utils.DurationToStr(startTime, time.Now()))

	return nil
}
