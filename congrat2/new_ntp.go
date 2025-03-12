package congrat2

import (
	"active/nts"
	"active/utils"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"
)

func executeNTP(ke *KeKeyTimestamp, aeadID int) error {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s: %d", ke.NTPv4Address, ke.NTPv4Port))
	if err != nil {
		return err
	}
	// 生成请求数据
	var req []byte
	var singleLen int
	if aeadID == 0 {
		req = utils.SecData()
	} else {
		// 消耗一个 Cookie
		eightLen := len(ke.Cookies)
		if eightLen == 0 {
			return errors.New("no cookie")
		}
		singleLen = eightLen / 8
		if singleLen*8 != eightLen {
			return fmt.Errorf("cookie length error: %d", eightLen)
		}
		req, err = nts.GenerateSecureNTPRequest(ke.C2SKey, ke.Cookies[:singleLen])
		ke.Cookies = ke.Cookies[singleLen:]
		if err != nil {
			return err
		}
	}
	// 建立连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	// 写数据
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	// 后面记得给 aeadID 为 0 的情况再赋值以免 T1!=T1R
	ke.T1R = utils.GetTimestamp(utils.GlobalNowTime())
	_, err = conn.Write(req)
	if err != nil {
		return err
	}

	// 接收响应
	buf := make([]byte, 1024)
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	if aeadID != 0 {
		// 验证响应
		cookieBuf := new(bytes.Buffer)
		err = nts.ValidateResponse(buf[:n], ke.S2CKey, cookieBuf)
		if err != nil {
			return err
		}
		ke.Cookies = append(ke.Cookies, cookieBuf.Bytes()[:singleLen]...)
		fmt.Printf("total cookies length after regain: %d\n", len(ke.Cookies))
	}
	// 记录 NTP 字段
	ke.Stratum = int(buf[1])
	ke.Poll = int(int8(buf[2]))
	ke.NTPPrecision = int(int8(buf[3]))
	ke.RootDelay = buf[4:8]
	ke.RootDispersion = buf[8:12]
	ke.Reference = buf[12:16]
	// 记录时间戳
	ke.T1, ke.T2, ke.T3 = buf[24:32], buf[32:40], buf[40:48]
	ke.T4 = utils.GetTimestamp(utils.GlobalNowTime())
	return nil
}

func queryPort(db *sql.DB, ke *KeKeyTimestamp) error {
	query := `SELECT ntpv4_address, ntpv4_port FROM ke_servers WHERE ip_address = ? ORDER BY id DESC LIMIT 1;`
	rows, err := db.Query(query, ke.IPAddress)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		addr, port := "", 0
		err = rows.Scan(&addr, &port)
		if err != nil {
			return err
		}
		ke.NTPv4Address, ke.NTPv4Port = addr, port
		return nil
	}
	return errors.New("no record found")
}
