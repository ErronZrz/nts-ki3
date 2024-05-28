package offset

import (
	"active/utils"
	"bytes"
	"net"
	"strings"
	"sync"
	"time"
)

func RecordNTPTimestamps(addr string, wg *sync.WaitGroup, errCh chan<- error) {
	defer func() {
		wg.Done()
	}()

	// 解析地址
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		// 检查是否是地址解析错误
		if !strings.Contains(err.Error(), "no such host") {
			errCh <- err
		}
		return
	}

	// 获取数据
	buf := new(bytes.Buffer)
	_, err = buf.Write(utils.SecData())
	if err != nil {
		errCh <- err
		return
	}

	// 发送数据
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		errCh <- err
		return
	}

	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Write(buf.Bytes())
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
	info := ServerTimestampsMap[0]
	info.T4 = time.Now()
	info.T1 = utils.ParseTimestamp(data[24:32])
	info.T2 = utils.ParseTimestamp(data[32:40])
	info.T3 = utils.ParseTimestamp(data[40:48])
}
