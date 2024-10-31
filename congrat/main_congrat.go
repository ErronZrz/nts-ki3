package congrat

import (
	"active/datastruct"
	"active/nts"
	"active/offset"
	"active/utils"
	"net"
	"sync"
	"time"
)

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
	var req []byte
	if aeadID == 0 {
		req = utils.SecData()
	} else {
		info.RLock()
		c2s, cookie := info.C2SKeyMap[aeadID], info.CookieMap[aeadID][0]
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
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		errCh <- err
		return
	}
	// 记录时间戳
	info.Lock()
	info.T4[aeadID] = utils.GlobalNowTime()
	info.Timestamps[aeadID] = buf[24:48]
	info.Unlock()
}
