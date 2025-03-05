package congrat

import (
	"active/datastruct"
	"active/nts"
	"active/offset"
	"active/utils"
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	CurrentBatchID int64
)

func MainFunction(path string, maxCoroutines int) error {
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

	var ipList []string

	for scanner.Scan() {
		ip := strings.Split(scanner.Text(), "\t")[0]
		ipList = append(ipList, ip)

		wg.Add(1)
		// 尝试向信号量发送数据，如果信号量满则会阻塞
		sem <- struct{}{}
		count++
		fmt.Printf("Start NTS-KE %d\n", count)

		go func(ip string) {
			ExecuteServerNTSKE(ip, errCh)
			wg.Done()
			<-sem // 释放信号量
			finished++
			fmt.Printf("Finished NTS-KE %d\n", finished)
		}(ip)
	}

	wg.Wait()
	count, finished = 0, 0

	for _, ip := range ipList {
		wg.Add(1)
		sem <- struct{}{}
		count++
		fmt.Printf("Start NTP %d\n", count)

		go func(ip string) {
			ExecuteServerSecureNTP(ip, errCh)
			wg.Done()
			<-sem // 释放信号量
			finished++
			fmt.Printf("Finished NTP %d\n", finished)
		}(ip)
	}

	wg.Wait()
	return nil
}

func ExecuteServerNTSKE(ip string, errCh chan<- error) {
	datastruct.OffsetMapMu.Lock()
	info, ok := datastruct.OffsetInfoMap[ip]
	if !ok {
		info = datastruct.NewOffsetServerInfo(ip)
		datastruct.OffsetInfoMap[ip] = info
	}
	datastruct.OffsetMapMu.Unlock()

	ipWg := new(sync.WaitGroup)
	ipWg.Add(3)

	go offset.AsyncRecordNTSTimestamps(ip, 0x0F, ipWg, errCh, true)
	go offset.AsyncRecordNTSTimestamps(ip, 0x10, ipWg, errCh, true)
	go offset.AsyncRecordNTSTimestamps(ip, 0x11, ipWg, errCh, true)

	ipWg.Wait()
}

func ExecuteServerSecureNTP(ip string, errCh chan<- error) {
	datastruct.OffsetMapMu.RLock()
	info := datastruct.OffsetInfoMap[ip]
	datastruct.OffsetMapMu.RUnlock()

	ipWg := new(sync.WaitGroup)
	ipWg.Add(4)

	go AsyncExecuteAEAD(0x00, ipWg, errCh, info)
	go AsyncExecuteAEAD(0x0F, ipWg, errCh, info)
	go AsyncExecuteAEAD(0x10, ipWg, errCh, info)
	go AsyncExecuteAEAD(0x11, ipWg, errCh, info)

	ipWg.Wait()
}

func AsyncExecuteAEAD(aeadID byte, wg *sync.WaitGroup, errCh chan<- error, info *datastruct.OffsetServerInfo) {
	defer wg.Done()

	// 如果不支持该算法则直接结束
	if aeadID != 0 {
		info.RLock()
		cookies := info.CookieMap[aeadID]
		info.RUnlock()
		if len(cookies) == 0 {
			return
		}
	}

	// 解析地址
	serverAddr := info.Server + ":" + info.Port
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		errCh <- err
		return
	}
	// 生成请求数据
	var req, s2c []byte
	if aeadID == 0 {
		req = utils.SecData()
	} else {
		info.RLock()
		c2s, cookie := info.C2SKeyMap[aeadID], info.CookieMap[aeadID][0]
		s2c = info.S2CKeyMap[aeadID]
		info.RUnlock()
		req, err = nts.GenerateSecureNTPRequest(c2s, cookie)
		if err != nil {
			errCh <- err
			return
		}
	}
	// 建立连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		errCh <- err
		return
	}
	defer func() { _ = conn.Close() }()
	// 写数据
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	info.Lock()
	// 这里可能会导致对于普通 NTP RealT1 与 T1 不同，但是后面补上了相关赋值，所以不影响
	info.RealT1[aeadID] = utils.GlobalNowTime()
	if aeadID != 0 {
		// 消耗一个 Cookie
		info.CookieMap[aeadID] = info.CookieMap[aeadID][1:]
		if len(info.CookieMap[aeadID]) == 0 {
			// 进行标记以便重新握手
			info.C2SKeyMap[aeadID] = nil
		}
	}
	info.Unlock()
	_, err = conn.Write(req)
	if err != nil {
		errCh <- err
		return
	}
	// 接收响应
	buf := make([]byte, 1024)
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		errCh <- err
		return
	}
	if aeadID != 0 {
		// 验证响应
		cookieBuf := new(bytes.Buffer)
		err = nts.ValidateResponse(buf[:n], s2c, cookieBuf)
		if err != nil {
			errCh <- err
			return
		}
		info.Lock()
		info.CookieMap[aeadID] = append(info.CookieMap[aeadID], cookieBuf.Bytes())
		info.Unlock()
	}
	info.Lock()
	info.PacketLen[aeadID] = n
	// 记录 NTP 字段
	info.Strata[aeadID] = int(buf[1])
	info.Polls[aeadID] = int(int8(buf[2]))
	info.Precisions[aeadID] = int(int8(buf[3]))
	info.RootDelays[aeadID] = buf[4:8]
	info.RootDispersions[aeadID] = buf[8:12]
	info.References[aeadID] = buf[12:16]
	// 记录时间戳
	info.T4[aeadID] = utils.GlobalNowTime()
	info.Timestamps[aeadID] = buf[24:48]
	info.Unlock()
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
