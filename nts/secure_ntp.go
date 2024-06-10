package nts

import (
	"active/utils"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"github.com/secure-io/siv-go"
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
	err = addAuthEF(buf, c2s)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

func addAuthEF(buf *bytes.Buffer, c2s []byte) error {
	algorithm, err := siv.NewCMAC(c2s)
	if err != nil {
		return err
	}

	nonce := make([]byte, 16)
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	cipherText := algorithm.Seal(nil, nonce, nil, buf.Bytes())

	nonceBuf := new(bytes.Buffer)
	err = binary.Write(nonceBuf, binary.BigEndian, nonce)
	if err != nil {
		return err
	}

	cipherBuf := new(bytes.Buffer)
	err = binary.Write(cipherBuf, binary.BigEndian, cipherText)
	if err != nil {
		return err
	}

	extensionBuf := new(bytes.Buffer)

	err = binary.Write(extensionBuf, binary.BigEndian, uint16(nonceBuf.Len()))
	if err != nil {
		return err
	}

	err = binary.Write(extensionBuf, binary.BigEndian, uint16(cipherBuf.Len()))
	if err != nil {
		return err
	}

	_, err = extensionBuf.ReadFrom(nonceBuf)
	if err != nil {
		return err
	}

	noncePadding := make([]byte, (nonceBuf.Len()+3) & ^3)
	_, err = extensionBuf.Write(noncePadding)
	if err != nil {
		return err
	}

	_, err = extensionBuf.ReadFrom(cipherBuf)
	if err != nil {
		return err
	}

	cipherPadding := make([]byte, (cipherBuf.Len()+3) & ^3)
	_, err = extensionBuf.Write(cipherPadding)
	if err != nil {
		return err
	}

	_, err = buf.Write([]byte{0x04, 0x04, 0x00, 0x28})
	if err != nil {
		return err
	}

	_, err = buf.ReadFrom(extensionBuf)
	return err
}
