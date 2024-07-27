package offset

import (
	"active/datastruct"
	"active/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func CalculateOffsetsAsync(inPath, outPath string, interval int) error {
	// 创建文件
	file, err := os.Open(inPath)
	if err != nil {
		return err
	}

	// 创建结果文件
	resultFile, err := os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
			if !errContains(err, "i/o timeout", "deadline exceeded", "failed to respond", "no such host") {
				fmt.Println(err)
			}
		}
	}()

	for scanner.Scan() {
		ip := strings.Split(scanner.Text(), "\t")[0]
		wg.Add(1)
		go func() {
			CalculateIPOffset(ip, errCh, 3)
			wg.Done()
		}()
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}

	wg.Wait()

	for ip, info := range datastruct.OffsetInfoMap {
		_, err = writer.WriteString(generateLine1(ip, info))
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func NTSKEDetectorWithFlags(path string, maxCoroutines int) error {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	wg := new(sync.WaitGroup)
	errCh := make(chan error, 16)
	// 创建大小为 maxCoroutines 的信号量
	sem := make(chan struct{}, maxCoroutines)
	var count, finished int

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
		count++
		fmt.Printf("Start %d\n", count)

		go func() {
			CalculateIPOffset(ip, errCh, 3)
			wg.Done()
			<-sem // 释放信号量
			finished++
			fmt.Printf("Finished %d\n", finished)
		}()
	}

	wg.Wait()
	return nil
}

func CalculateIPOffset(ip string, errCh chan<- error, aeadNum int) {
	datastruct.OffsetMapMu.Lock()
	info, ok := datastruct.OffsetInfoMap[ip]
	if !ok {
		info = datastruct.NewOffsetServerInfo(ip)
		datastruct.OffsetInfoMap[ip] = info
	}
	datastruct.OffsetMapMu.Unlock()

	ipWg := new(sync.WaitGroup)
	ipWg.Add(3)

	go AsyncRecordNTSTimestamps(ip, 0x0F, ipWg, errCh, aeadNum < 1)
	go AsyncRecordNTSTimestamps(ip, 0x10, ipWg, errCh, aeadNum < 2)
	go AsyncRecordNTSTimestamps(ip, 0x11, ipWg, errCh, aeadNum < 3)

	ipWg.Wait()

	if aeadNum > 1 {
		AsyncRecordNTPTimestamps(info, errCh)
	}
}

func generateLine1(ip string, info *datastruct.OffsetServerInfo) string {
	var classStr string
	addStr := func(b bool) {
		if b {
			classStr += "Y"
		} else {
			classStr += "N"
		}
	}
	addStr(info.RightIP)
	addStr(!info.Expired)
	addStr(len(info.C2SKeyMap[0x0F]) > 0)
	addStr(len(info.CookieMap[0x0F]) > 0)
	addStr(!info.T1[0x0F].IsZero())
	line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		ip, info.CommonName, classStr,
		getOffset1(info, 0x00, false),
		getOffset1(info, 0x0F, false),
		getOffset1(info, 0x0F, true),
		getOffset1(info, 0x10, false),
		getOffset1(info, 0x10, true),
		getOffset1(info, 0x11, false),
		getOffset1(info, 0x11, true),
	)
	info.ClearTimeStamps()
	return line
}

func GenerateLine2(ip string, info *datastruct.OffsetServerInfo) string {
	cookieMap := info.CookieMap
	cookies256 := cookieMap[0x0F]
	if len(cookies256) == 0 {
		return ""
	}

	server := info.Server
	if server == ip {
		server = "Default"
	}

	aeadStr := fmt.Sprintf("SIV_CMAC_256(%d)", len(cookies256[0]))
	if len(cookieMap[0x10]) > 0 {
		aeadStr += fmt.Sprintf(",SIV_CMAC_384(%d)", len(cookieMap[0x10][0]))
	}
	if len(cookieMap[0x11]) > 0 {
		aeadStr += fmt.Sprintf(",SIV_CMAC_512(%d)", len(cookieMap[0x11][0]))
	}

	var flagStr string
	addStr := func(b bool) {
		if b {
			flagStr += "Y"
		} else {
			flagStr += "N"
		}
	}
	addStr(info.RightIP)
	addStr(!info.Expired)
	addStr(!info.SelfSigned)
	addStr(!info.T1[0x0F].IsZero())

	var keyIDStr string
	if len(cookies256) >= 2 {
		keyIDStr = utils.SameFourBytes(cookies256[0], cookies256[1])
	}
	if len(keyIDStr) == 0 {
		keyIDStr = "no-keyID"
	}

	dateFormat := "2006-01-02 15:04:05"
	strList := []string{
		ip,
		info.CommonName,
		server,
		info.Port,
		aeadStr,
		flagStr,
		keyIDStr,
		info.NotBefore.Format(dateFormat),
		info.NotAfter.Format(dateFormat),
		info.Organization,
		info.Issuer,
		time.Now().Format(dateFormat),
		getOffset1(info, 0x00, false),
		getOffset1(info, 0x0F, false),
		getOffset1(info, 0x0F, true),
		getOffset1(info, 0x10, false),
		getOffset1(info, 0x10, true),
		getOffset1(info, 0x11, false),
		getOffset1(info, 0x11, true),
	}
	return strings.Join(strList, "\t") + "\n"
}

func getOffset1(info *datastruct.OffsetServerInfo, aeadID byte, useReal bool) string {
	t1, t2, t3, t4 := info.T1[aeadID], info.T2[aeadID], info.T3[aeadID], info.T4[aeadID]
	if useReal {
		t1 = info.RealT1[aeadID]
	}
	if t2.IsZero() || t1.IsZero() || t3.IsZero() || t4.IsZero() {
		return "-"
	}
	offset := (t2.Sub(t1.UTC()) + t3.Sub(t4.UTC())) / 2
	return fmt.Sprintf("%.3f", float64(offset.Nanoseconds()/1000)/1000)
}

func errContains(err error, substrList ...string) bool {
	errStr := err.Error()
	for _, substr := range substrList {
		if strings.Contains(errStr, substr) {
			return true
		}
	}
	return false
}
