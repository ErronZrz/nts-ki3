package nts

import (
	"active/parser"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
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

	// 新建 reader 和 buf
	reader := bufio.NewReader(file)

	// 遍历每行
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// 读取第一个制表符的下标
		i := bytes.IndexByte(line, '\t')

		// 执行探测并解析
		ip := line[:i]
		payload, err := DialNTSKE(string(ip), "", 0x22)
		if err != nil {
			return err
		}
		if payload.Len == 0 {
			return errors.New("NTS-KE payload is empty")
		}
		_, err = parser.ParseNTSResponse(payload.RcvData)
		if err != nil {
			return err
		}

		// 将服务器地址写入文件
		if parser.TheServer == "" {
			parser.TheServer = string(ip)
		}
		_, err = outputFile.WriteString(parser.TheServer + " ")
		parser.TheServer = ""
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
