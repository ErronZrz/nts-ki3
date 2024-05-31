package offset

import (
	"active/nts"
	"active/parser"
	"active/utils"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type ServerTimestamps struct {
	T1     time.Time
	RealT1 time.Time
	T2     time.Time
	T3     time.Time
	T4     time.Time
}

var (
	ServerTimestampsMap map[byte]*ServerTimestamps
)

func init() {
	ServerTimestampsMap = map[byte]*ServerTimestamps{
		0x00: new(ServerTimestamps),
		0x0F: new(ServerTimestamps),
		0x10: new(ServerTimestamps),
		0x11: new(ServerTimestamps),
	}
}

func RecordNTSTimestamps(ip string, aeadID byte, wg *sync.WaitGroup, errCh chan<- error) {
	defer func() {
		wg.Done()
	}()
	// 执行 NTS 握手
	payload, err := nts.DialNTSKE(ip, "", aeadID)
	if err != nil {
		// 检查是否是超时错误
		errStr := err.Error()
		if !strings.Contains(errStr, "i/o timeout") && !strings.Contains(errStr, "deadline exceeded") {
			errCh <- err
		}
		return
	}
	// 记录 C2S 密钥
	c2sKey := payload.C2SKey
	if payload.Len == 0 {
		errCh <- errors.New(fmt.Sprintf("%s: NTS-KE payload is empty", ip))
		return
	}
	// 解析并记录 Cookie
	_, err = parser.ParseNTSResponse(payload.RcvData)
	if err != nil {
		errCh <- err
		return
	}
	// 如果没有 Cookie 说明不支持该算法，直接结束
	parser.MuCookie.RLock()
	cookie := parser.CookieMap[aeadID]
	parser.MuCookie.RUnlock()
	if cookie == nil {
		return
	}
	// 记录服务器地址和端口号
	serverAddr := parser.TheHost
	if serverAddr == "" {
		serverAddr = ip
	}
	serverAddr += ":" + parser.ThePort

	// 解析地址
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		// 检查是否是地址解析错误
		if !strings.Contains(err.Error(), "no such host") {
			errCh <- err
		}
		return
	}
	// 生成请求数据
	req, err := nts.GenerateSecureNTPRequest(c2sKey, cookie)
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
	info := ServerTimestampsMap[aeadID]
	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	info.RealT1 = time.Now()
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
	info.T4 = time.Now()
	info.T1 = utils.ParseTimestamp(buf[24:32])
	info.T2 = utils.ParseTimestamp(buf[32:40])
	info.T3 = utils.ParseTimestamp(buf[40:48])
}
