package nts

import (
	"active/utils"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"github.com/secure-io/siv-go"
)

var (
	PlaceholderNum int
)

func GenerateSecureNTPRequest(c2s, cookie []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := addHeader(buf)
	if err != nil {
		return nil, err
	}
	err = addUniqueEF(buf)
	if err != nil {
		return nil, err
	}
	err = addCookieEF(buf, cookie)
	if err != nil {
		return nil, err
	}
	if PlaceholderNum > 0 && PlaceholderNum < 8 {
		cookieSize := len(cookie)
		for i := 0; i < PlaceholderNum; i++ {
			err = addCookiePlaceholderEF(buf, cookieSize)
			if err != nil {
				return nil, err
			}
		}
	}
	err = addAuthEF(buf, c2s)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ValidateResponse 解析并验证 NTS 服务器响应
func ValidateResponse(data, key []byte, cookieBuf *bytes.Buffer) error {
	if len(data) < 160 {
		return errors.New("data length is too short")
	}
	ad := make([]byte, 84)
	copy(ad, data)
	cipherLen := 256*int(data[90]) + int(data[91])
	nonce := make([]byte, 16)
	copy(nonce, data[92:])
	cipherText := make([]byte, cipherLen)
	copy(cipherText, data[108:])

	// 创建 SIV-CMAC AEAD 实例
	aead, err := siv.NewCMAC(key)
	if err != nil {
		return err
	}

	// 使用提供的 nonce, cipherText 和 ad 解密并验证数据
	// 该过程已经通过比较摘要验证了数据完整性，不需要手动检查
	result, err := aead.Open(nil, nonce, cipherText, ad)
	if err != nil {
		return err
	}

	// 将解密后的 Cookie 数据写入 buffer
	cookieBuf.Write(result)

	return nil
}

func addHeader(buf *bytes.Buffer) error {
	_, err := buf.Write(utils.SecData())
	return err
}

func addUniqueEF(buf *bytes.Buffer) error {
	_, err := buf.Write([]byte{0x01, 0x04, 0x00, 0x24})
	if err != nil {
		return err
	}
	// 添加随机的 32 字节
	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return err
	}
	_, err = buf.Write(randomBytes)
	return err
}

func addCookieEF(buf *bytes.Buffer, cookie []byte) error {
	_, err := buf.Write([]byte{0x02, 0x04, 0x00, byte(len(cookie)) + 4})
	if err != nil {
		return err
	}
	_, err = buf.Write(cookie)
	return err
}

func addCookiePlaceholderEF(buf *bytes.Buffer, cookieSize int) error {
	_, err := buf.Write([]byte{0x03, 0x04, 0x00, byte(cookieSize) + 4})
	if err != nil {
		return err
	}
	emptyCookie := make([]byte, cookieSize)
	_, err = buf.Write(emptyCookie)
	return err
}

// addAuthEF 函数用于向提供的 buffer 中添加一个认证扩展字段
func addAuthEF(buf *bytes.Buffer, c2s []byte) error {
	// 初始化 CMAC 算法
	algorithm, err := siv.NewCMAC(c2s)
	if err != nil {
		return err
	}

	// 生成一个随机 nonce（16 字节）
	nonce := make([]byte, 16)
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	// 使用 CMAC 算法加密 buf 中的数据
	cipherText := algorithm.Seal(nil, nonce, nil, buf.Bytes())

	// 将 nonce 写入一个新的 buffer，使用大端格式
	nonceBuf := new(bytes.Buffer)
	err = binary.Write(nonceBuf, binary.BigEndian, nonce)
	if err != nil {
		return err
	}

	// 将加密后的文本写入另一个新的 buffer，同样使用大端格式
	cipherBuf := new(bytes.Buffer)
	err = binary.Write(cipherBuf, binary.BigEndian, cipherText)
	if err != nil {
		return err
	}

	// 创建一个用于存放扩展字段的 buffer
	extensionBuf := new(bytes.Buffer)

	// 写入 nonceBuf 和 cipherBuf 的长度信息
	err = binary.Write(extensionBuf, binary.BigEndian, uint16(nonceBuf.Len()))
	if err != nil {
		return err
	}
	err = binary.Write(extensionBuf, binary.BigEndian, uint16(cipherBuf.Len()))
	if err != nil {
		return err
	}

	// 将 nonceBuf 的内容读取到 extensionBuf 中
	_, err = extensionBuf.ReadFrom(nonceBuf)
	if err != nil {
		return err
	}

	// 为 nonce 添加必要的填充以确保对齐
	noncePadding := make([]byte, (nonceBuf.Len()+3) & ^3)
	_, err = extensionBuf.Write(noncePadding)
	if err != nil {
		return err
	}

	// 将 cipherBuf 的内容读取到 extensionBuf 中
	_, err = extensionBuf.ReadFrom(cipherBuf)
	if err != nil {
		return err
	}

	// 为加密文本添加必要的填充以确保对齐
	cipherPadding := make([]byte, (cipherBuf.Len()+3) & ^3)
	_, err = extensionBuf.Write(cipherPadding)
	if err != nil {
		return err
	}

	// 将一个预定义的字节序列写入 buf
	_, err = buf.Write([]byte{0x04, 0x04, 0x00, 0x28})
	if err != nil {
		return err
	}

	// 将 extensionBuf 的内容添加到 buf 中
	_, err = buf.ReadFrom(extensionBuf)
	return err
}
