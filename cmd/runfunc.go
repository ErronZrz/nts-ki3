package cmd

import (
	"active/async"
	"active/datastruct"
	"active/output"
	"active/parser"
	"active/udpdetect"
	"active/utils"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

const (
	npLimit = 16
)

func executeTimeSync(cmd *cobra.Command, args []string) error {
	if nPrintedHosts > npLimit {
		nPrintedHosts = npLimit
	}
	cmdName := cmd.Name()
	if args == nil || len(args) == 0 {
		return errors.New("command `timesync` missing arguments")
	}
	address := args[0]
	var ngStr string
	if nGoroutines <= 0 {
		ngStr = "auto"
	} else {
		ngStr = strconv.Itoa(nGoroutines)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Ready to run `%s`.\n    address: %s\n    num of goroutines: %s\n"+
		"    num of printed hosts: %d\n\n", cmdName, address, ngStr, nPrintedHosts)

	var dataCh <-chan *datastruct.RcvPayload
	startTime := time.Now()
	if nGoroutines <= 0 {
		dataCh = udpdetect.DialNetworkNTP(address)
	} else {
		dataCh = udpdetect.DialNetworkNTPWithBatchSize(address, nGoroutines)
	}
	if dataCh == nil {
		return errors.New("dataCh is nil")
	}
	count := printResult(dataCh, "timesync_"+address)

	_, _ = fmt.Fprintf(os.Stdout, "%d hosts detected in %s\n",
		count, utils.DurationToStr(startTime, time.Now()))

	return nil
}

func executeAsync(cmd *cobra.Command, args []string) error {
	if nPrintedHosts > npLimit {
		nPrintedHosts = npLimit
	}
	cmdName := cmd.Name()
	if args == nil || len(args) == 0 {
		return errors.New("command `async` missing arguments")
	}
	address := args[0]
	_, _ = fmt.Fprintf(os.Stdout, "Ready to run `%s`.\n    address: %s\n    "+
		"num of printed hosts: %d\n\n", cmdName, address, nPrintedHosts)

	startTime := time.Now()
	dataCh := async.DialNetworkNTP(address, 1)

	if dataCh == nil {
		return errors.New("dataCh is nil")
	}

	count := printResult(dataCh, "async_"+address)

	_, _ = fmt.Fprintf(os.Stdout, "%d hosts detected in %s\n",
		count, utils.DurationToStr(startTime, time.Now()))

	return nil
}

func printResult(dataCh <-chan *datastruct.RcvPayload, cmd string) int {
	seqNum := 0
	now := time.Now()

	for p, ok := <-dataCh; ok; p, ok = <-dataCh {
		err := p.Err
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			continue
		}
		header, err := parser.ParseHeader(p.RcvData)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
		} else {
			seqNum++
			payloadStr, headerStr := p.Lines(), header.Lines()
			output.WriteToFile(payloadStr, headerStr, cmd, seqNum, p.RcvTime, now)
			if seqNum <= nPrintedHosts {
				_, _ = fmt.Fprintf(os.Stdout, "[Host %d]\n", seqNum)
				_, _ = fmt.Fprint(os.Stdout, payloadStr)
				_, _ = fmt.Fprintln(os.Stdout, "[parsed]")
				_, _ = fmt.Fprint(os.Stdout, headerStr)
			}
		}
	}

	return seqNum
}

//TODO: 优化 WriteToFile，创建文件时先输出本次扫描相关信息
//TODO: 地址随机化
