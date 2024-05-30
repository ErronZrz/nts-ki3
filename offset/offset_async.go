package offset

import (
	"active/datastruct"
	"active/nts"
	"active/parser"
	"active/utils"
	"errors"
	"net"
	"strings"
	"sync"
	"time"
)

func AsyncRecordNTSTimestamps(ip string, aeadID byte, wg *sync.WaitGroup, errCh chan<- error) {
	defer func() {
		wg.Done()
	}()

	datastruct.OffsetMapMu.RLock()
	info := datastruct.OffsetInfoMap[ip]
	datastruct.OffsetMapMu.RUnlock()

	// 如果没有进行过 KE 握手则进行
	info.RLock()
	c2s := info.C2SKeyMap[aeadID]
	info.RUnlock()
	if c2s == nil {
		// 标记为已经进行过
		info.Lock()
		info.C2SKeyMap[aeadID] = []byte{}
		info.Unlock()

		payload, err := nts.DialNTSKE(ip, "", aeadID)
		if err != nil {
			// 检查是否是超时错误
			errStr := err.Error()
			if !strings.Contains(errStr, "i/o timeout") && !strings.Contains(errStr, "deadline exceeded") {
				errCh <- err
			}
			return
		}
		if payload.Len == 0 {
			errCh <- errors.New("NTS-KE payload is empty")
			return
		}
		info.Lock()
		info.C2SKeyMap[aeadID] = payload.C2SKey
		info.Unlock()

		// 解析 Cookie 等信息
		err = parser.ParseOffsetInfo(payload.RcvData, info, aeadID)
		if err != nil {
			errCh <- err
			return
		}
	}

	// 如果不支持该算法则直接结束
	info.RLock()
	cookies := info.CookieMap[aeadID]
	info.RUnlock()
	if len(cookies) == 0 {
		return
	}

	// 解析地址
	serverAddr := info.Server + ":" + info.Port
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		// 检查是否是地址解析错误
		if !strings.Contains(err.Error(), "no such host") {
			errCh <- err
		}
		return
	}
	// 生成请求数据
	info.RLock()
	c2s, cookie := info.C2SKeyMap[aeadID], info.CookieMap[aeadID][0]
	info.RUnlock()
	req, err := nts.GenerateSecureNTPRequest(c2s, cookie)
	if err != nil {
		errCh <- err
		return
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
	info.RealT1[aeadID] = time.Now()
	info.CookieMap[aeadID] = info.CookieMap[aeadID][1:]
	info.Unlock()
	_, err = conn.Write(req)
	if err != nil {
		errCh <- err
		return
	}
	// 接收响应
	buf := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		// 检查是否是超时错误
		if !strings.Contains(err.Error(), "i/o timeout") {
			errCh <- err
		}
		return
	}
	// 记录时间戳
	info.Lock()
	info.T4[aeadID] = time.Now()
	info.T1[aeadID] = utils.ParseTimestamp(buf[24:32])
	info.T2[aeadID] = utils.ParseTimestamp(buf[32:40])
	info.T3[aeadID] = utils.ParseTimestamp(buf[40:48])
	info.Unlock()
}

func AsyncRecordNTPTimestamps(info *datastruct.OffsetServerInfo, errCh chan<- error) {
	// 解析地址
	serverAddr := info.Server + ":" + info.Port
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		// 检查是否是地址解析错误
		if !strings.Contains(err.Error(), "no such host") {
			errCh <- err
		}
		return
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
	_, err = conn.Write(utils.SecData())
	if err != nil {
		errCh <- err
		return
	}
	// 接收响应
	data := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(data)
	if err != nil {
		// 检查是否是超时错误
		if !strings.Contains(err.Error(), "i/o timeout") {
			errCh <- err
		}
		return
	}
	// 记录时间戳
	info.Lock()
	info.T4[0x00] = time.Now()
	info.T1[0x00] = utils.ParseTimestamp(data[24:32])
	info.T2[0x00] = utils.ParseTimestamp(data[32:40])
	info.T3[0x00] = utils.ParseTimestamp(data[40:48])
	info.Unlock()
}
