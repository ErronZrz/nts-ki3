package nts

import (
	"active/parser"
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func CollectKeys(path string) error {
	// 打开 TXT 文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// 创建输出文件
	outputFilePath := path[:len(path)-4] + "_keys.txt"
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)

	// 遍历每行
	for scanner.Scan() {
		line := scanner.Text()
		// 读取第一个制表符的下标
		i := strings.IndexByte(line, '\t')

		// 执行探测并解析
		ip := line[:i]
		payload, err := DialNTSKE(ip, "", 0x11)
		if err != nil {
			// 检查是否是超时错误
			errStr := err.Error()
			if strings.Contains(errStr, "i/o timeout") || strings.Contains(errStr, "deadline exceeded") {
				fmt.Println(ip + " " + errStr)
				continue
			}
			return err
		}
		if payload.Len == 0 {
			return errors.New(fmt.Sprintf("%s: NTS-KE payload is empty", ip))
		}
		_, err = parser.ParseNTSResponse(payload.RcvData)
		if err != nil {
			return err
		}

		// 将服务器地址写入文件
		if parser.TheHost == "" {
			parser.TheHost = ip
		}
		_, err = outputFile.WriteString(parser.TheHost + ":" + parser.ThePort + " ")
		parser.TheHost, parser.ThePort = "", "123"
		if err != nil {
			return err
		}
		// 将密钥写入文件
		for _, b := range payload.C2SKey {
			_, err = outputFile.WriteString(fmt.Sprintf("%02X", b))
			if err != nil {
				return err
			}
		}
		_, err = outputFile.Write([]byte{' '})
		if err != nil {
			return err
		}
		// 将 Cookie 写入文件
		for _, b := range parser.FirstCookie {
			_, err = outputFile.WriteString(fmt.Sprintf("%02X", b))
			if err != nil {
				return err
			}
		}
		parser.FirstCookie = nil
		_, err = outputFile.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}

	// 完成
	return nil
}

func MakeSecureNTPRequests(path string) error {
	// 打开 TXT 文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)

	// 遍历每行
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		addr, keyStr, cookieStr := parts[0], parts[1], parts[2]

		// 解析密钥和 Cookie
		key, err := hex.DecodeString(keyStr)
		if err != nil {
			return err
		}
		cookie, err := hex.DecodeString(cookieStr)
		if err != nil {
			return err
		}

		// 生成请求数据
		req, err := GenerateSecureNTPRequest(key, cookie)
		if err != nil {
			return err
		}

		// 发送 UDP 数据
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			// 检查是否是地址解析错误
			if strings.Contains(err.Error(), "no such host") {
				fmt.Println(addr + " no such host")
				continue
			}
			return err
		}

		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			return err
		}
		_, err = conn.Write(req)
		if err != nil {
			return err
		}

		// 接收响应
		buf := make([]byte, 1024)
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		_, _, err = conn.ReadFromUDP(buf)
		if err != nil {
			// 检查是否是超时错误
			if strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println(addr + " timeout")
			} else {
				return err
			}
		}

		// 关闭连接
		err = conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
