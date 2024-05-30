package offset

import (
	"active/datastruct"
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func CalculateOffsetsAsync(inPath, outPath string) error {
	// 创建文件
	file, err := os.Open(inPath)
	if err != nil {
		return err
	}

	// 创建结果文件
	resultFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
		_ = resultFile.Close()
	}()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(resultFile)
	wg := new(sync.WaitGroup)
	errCh := make(chan error, 16)

	defer func() { close(errCh) }()

	go func() {
		for err := range errCh {
			fmt.Println(err)
		}
	}()

	for scanner.Scan() {
		ip := strings.Split(scanner.Text(), "\t")[0]
		wg.Add(1)
		go CalculateIPOffset(ip, wg, errCh)
		// 等待一秒
		time.Sleep(time.Second)
	}

	wg.Wait()

	for ip, info := range datastruct.OffsetInfoMap {
		if info.T1[0x0F].IsZero() {
			continue
		}
		_, err = writer.WriteString(generateLine1(ip, info))
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func CalculateIPOffset(ip string, wg *sync.WaitGroup, errCh chan<- error) {
	defer func() { wg.Done() }()
	info := datastruct.NewOffsetServerInfo(ip)
	datastruct.OffsetMapMu.Lock()
	datastruct.OffsetInfoMap[ip] = info
	datastruct.OffsetMapMu.Unlock()

	ipWg := new(sync.WaitGroup)
	ipWg.Add(3)

	go AsyncRecordNTSTimestamps(ip, 0x0F, ipWg, errCh)
	go AsyncRecordNTSTimestamps(ip, 0x10, ipWg, errCh)
	go AsyncRecordNTSTimestamps(ip, 0x11, ipWg, errCh)

	ipWg.Wait()

	AsyncRecordNTPTimestamps(info, errCh)
}

func generateLine1(ip string, info *datastruct.OffsetServerInfo) string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		ip,
		getOffset1(info, 0x00, false),
		getOffset1(info, 0x0F, false),
		getOffset1(info, 0x0F, true),
		getOffset1(info, 0x10, false),
		getOffset1(info, 0x10, true),
		getOffset1(info, 0x11, false),
		getOffset1(info, 0x11, true),
	)
}

func getOffset1(info *datastruct.OffsetServerInfo, aeadID byte, useReal bool) string {
	t1, t2, t3, t4 := info.T1[aeadID], info.T2[aeadID], info.T3[aeadID], info.T4[aeadID]
	if useReal {
		t1 = info.RealT1[aeadID]
	}
	if t1.IsZero() || t2.IsZero() || t3.IsZero() || t4.IsZero() {
		return "-"
	}
	offset := (t2.Sub(t1.UTC()) + t3.Sub(t4.UTC())) / 2
	return fmt.Sprintf("%.3f", float64(offset.Nanoseconds()/1000)/1000)
}
