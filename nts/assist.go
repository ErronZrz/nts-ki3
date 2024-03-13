package nts

import (
	"active/datastruct"
	"active/parser"
	"crypto/tls"
	"fmt"
	"io"
)

var (
	variableReq = []byte{
		0x80, 0x01, 0x00, 0x02, 0x00, 0x00, 0x80, 0x04, 0x00, 0x02, 0x00, 0x0F, 0x80, 0x00, 0x00, 0x00,
	}
	otherThanAesSivCmac = []byte{
		0x80, 0x01, 0x00, 0x02, 0x00, 0x00, 0x80, 0x04, 0x00, 0x3C, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03,
		0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x08, 0x00, 0x09, 0x00, 0x0A, 0x00, 0x0B,
		0x00, 0x0C, 0x00, 0x0D, 0x00, 0x0E, 0x00, 0x12, 0x00, 0x13, 0x00, 0x14, 0x00, 0x15, 0x00, 0x16,
		0x00, 0x17, 0x00, 0x18, 0x00, 0x19, 0x00, 0x1A, 0x00, 0x1B, 0x00, 0x1C, 0x00, 0x1D, 0x00, 0x1E,
		0x00, 0x1F, 0x00, 0x20, 0x00, 0x21, 0x80, 0x00, 0x00, 0x00,
	}
)

func singleReadWrite(aeadID byte, conn *tls.Conn, info *datastruct.DetectInfo) error {
	//defer func(conn *tls.Conn) {
	//	err := conn.Close()
	//	if err != nil {
	//		fmt.Printf("error closing TLS connection %v", err)
	//	}
	//}(conn)

	variableReq[11] = aeadID

	_, err := conn.Write(variableReq)
	if err != nil {
		return fmt.Errorf("send NTS-KE request failed: %v", err)
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("read NTS-KE response failed: %v", err)
	}
	_, _ = parser.ParseNTSResponse(data)
	/*
		if err != nil {
			fmt.Println("error")
		} else {
			fmt.Print(res.Lines())
		}*/
	return parser.ParseDetectInfo(data, info)
}

func checkOtherThanAesSivCmac(conn *tls.Conn, info *datastruct.DetectInfo) (bool, error) {
	defer func(conn *tls.Conn) {
		/*
			err := conn.Close()
			if err != nil {
				fmt.Printf("error closing TLS connection: %v", err)
			}*/
		_ = conn.Close()
	}(conn)

	_, err := conn.Write(otherThanAesSivCmac)
	if err != nil {
		return false, fmt.Errorf("send NTS-KE request failed: %v", err)
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		return false, fmt.Errorf("read NTS-KE response failed: %v", err)
	}

	err = parser.ParseDetectInfo(data, info)
	if err != nil {
		return false, err
	}

	support := false
	for id := byte(0x01); id <= 0x21; id++ {
		if id == 0x0F || id == 0x10 || id == 0x11 {
			continue
		}
		if info.AEADList[id] {
			support = true
			break
		}
	}
	return support, nil
}
