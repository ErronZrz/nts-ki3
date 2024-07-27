package offset

import (
	"active/nts"
	"active/parser"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	keyIDMap = make(map[string]string)
	mu       sync.Mutex
)

func GetNTSKeyID(inPath, outPath string, maxCoroutines int) error {
	// 打开文件
	file, err := os.Open(inPath)
	if err != nil {
		return err
	}
	outFile, err := os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
		_ = outFile.Close()
	}()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outFile)
	wg := new(sync.WaitGroup)
	errCh := make(chan error, 16)
	// 创建大小为 maxCoroutines 的信号量
	sem := make(chan struct{}, maxCoroutines)

	go func() {
		for err := range errCh {
			if !errContains(err, "i/o timeout", "deadline exceeded", "failed to respond", "no such host") {
				fmt.Println(err)
			}
		}
	}()

	for scanner.Scan() {
		ip := strings.Split(scanner.Text(), "\t")[0]

		wg.Add(1)
		// 尝试向信号量发送数据，如果信号量满则会阻塞
		sem <- struct{}{}

		go func() {
			getIPKeyID(ip, errCh)
			wg.Done()
			<-sem
		}()
	}

	wg.Wait()

	for ip, keyID := range keyIDMap {
		_, err := writer.WriteString(ip + "\t" + keyID + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func getIPKeyID(ip string, errCh chan error) {
	payload, err := nts.DialNTSKE(ip, "", 0x0F)
	if err != nil {
		errCh <- err
		return
	}

	if payload.Len == 0 {
		errCh <- errors.New(fmt.Sprintf("%s: NTS-KE payload is empty", ip))
		return
	}

	keyID, err := parser.ParseKeyID(payload.RcvData)
	if err != nil {
		errCh <- err
		return
	}

	mu.Lock()
	keyIDMap[ip] = keyID
	mu.Unlock()
}
