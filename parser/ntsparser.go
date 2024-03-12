package parser

import (
	"active/datastruct"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

const (
	defaultNTPPortStr = "123"
)

var (
	recordTypeNames = []string{
		"End of Message", "Next Protocol", "Error", "Warning",
		"AEAD Algorithm", "Cookie", "NTPv4 Server", "NTPv4 Port",
	}
	printers = []printer{
		printEOM, printNextProto, printError, printWarning,
		printAEADAlgorithm, printCookie, printServer, printPort,
	}
)
var (
	reqBytes = []byte{
		0x80, 0x01, 0x00, 0x02, 0x00, 0x00, 0x80, 0x04, 0x00, 0x02, 0x00, 0x0F, 0x80, 0x00, 0x00, 0x00,
	}
)

type Response struct {
	buf *bytes.Buffer
}

type printer func(*record, *bytes.Buffer) error

type record struct {
	critical bool
	rType    uint16
	bodyLen  uint16
	body     []byte
}

func ParseNTSResponse(data []byte) (*Response, error) {
	records, err := retrieveRecords(data)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	for i, next := range records {
		var lenStr, criticalStr, newLineStr string
		if next.bodyLen == 0 {
			lenStr = "empty body"
		} else {
			lenStr = strconv.Itoa(int(next.bodyLen)) + "B"
		}
		if next.critical {
			criticalStr = " (C)"
		}
		if next.rType > 0 {
			newLineStr = "\n            "
		}
		typeName := recordTypeNames[next.rType]
		buf.WriteString(fmt.Sprintf("Record %d: %s (%s)%s%s", i+1, typeName, lenStr, criticalStr, newLineStr))
		err = printers[next.rType](next, buf)
		if err != nil {
			return nil, err
		}
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')

	return &Response{buf: buf}, nil
}

func (r *Response) Lines() string {
	return r.buf.String()
}

func ParseDetectInfo(data []byte, info datastruct.DetectInfo) error {
	records, err := retrieveRecords(data)
	if err != nil {
		return err
	}
	var server, portStr string
	for _, r := range records {
		switch r.rType {
		// AEAD Algorithm
		case 4:
			if r.bodyLen > 0 {
				if r.bodyLen != 2 {
					return fmt.Errorf("unexpected body length in Warning record: %d", r.bodyLen)
				}
				if r.body[0] != 0x00 || r.body[1] == 0x00 || r.body[1] > 0x21 {
					return fmt.Errorf("unrecognized AEAD algorithm ID in AEAD Algorhithm record: %d",
						(int(r.body[0])<<8)+int(r.body[1]))
				}
				info.AEADList[r.body[1]] = true
			}
		// Cookie
		case 5:
			if r.bodyLen == 0 {
				return errors.New("empty body in Cookie record")
			}
			info.CookieLength = int(r.bodyLen)
		// NTPv4 Server
		case 6:
			if r.bodyLen == 0 {
				return errors.New("empty body in NTPv4 Server record")
			}
			server = string(r.body)
		// NTPv4 Port
		case 7:
			if r.bodyLen != 2 {
				return fmt.Errorf("unexpected body length in NTPv4 Port record: %d", r.bodyLen)
			}
			port := (int(r.body[0]) << 8) + int(r.body[1])
			portStr = strconv.Itoa(port)
		default:
			continue
		}
	}

	if server != "" {
		if portStr == "" {
			portStr = defaultNTPPortStr
		}
		info.ServerPortSet[server+":"+portStr] = struct{}{}
	}

	return nil
}

func retrieveRecords(data []byte) ([]*record, error) {
	n := len(data)
	cur := 0
	res := make([]*record, 0)
	for cur < n {
		if cur+4 > n {
			return nil, errors.New("an incomplete record exists at the end")
		}
		next := new(record)
		if data[cur] == 0x80 {
			next.critical = true
		} else if data[cur] == 0x00 {
			next.critical = false
		} else {
			return nil, fmt.Errorf("invalid header byte %02X at pos %d", data[cur], cur)
		}
		cur++
		if data[cur] >= 8 {
			return nil, fmt.Errorf("invalid header byte %02X at pos %d", data[cur], cur)
		}
		next.rType = uint16(data[cur])
		cur++
		bodyLen := binary.BigEndian.Uint16(data[cur : cur+2])
		cur += 2
		if bodyLen == 0 {
			res = append(res, next)
			continue
		}
		endPos := cur + int(bodyLen)
		if endPos > n {
			return nil, errors.New("a record with an incorrect body length exists at the end")
		}
		next.bodyLen = bodyLen
		next.body = data[cur:endPos]
		cur = endPos
		res = append(res, next)
	}
	return res, nil
}

func printEOM(r *record, _ *bytes.Buffer) error {
	if r.bodyLen > 0 {
		return fmt.Errorf("unexpected body length in EOM record: %d", r.bodyLen)
	}
	return nil
}

func printNextProto(r *record, buf *bytes.Buffer) error {
	if r.bodyLen == 0 {
		buf.WriteString("None of given protocol supported")
		return nil
	}
	if r.bodyLen != 2 {
		return fmt.Errorf("unexpected body length in Next Protocol record: %d", r.bodyLen)
	}
	if r.body[0] != 0x00 || r.body[1] != 0x00 {
		return fmt.Errorf("unrecognizable protocol ID in Next Protocol record: %d",
			(int(r.body[0])<<8)+int(r.body[1]))
	}
	buf.WriteString("Protocol = NTPv4")
	return nil
}

func printError(r *record, buf *bytes.Buffer) error {
	if r.bodyLen != 2 {
		return fmt.Errorf("unexpected body length in Error record: %d", r.bodyLen)
	}
	if r.body[0] != 0x00 || r.body[1] > 0x02 {
		return fmt.Errorf("unrecognized error code in Error record: %d",
			(int(r.body[0])<<8)+int(r.body[1]))
	}
	var explain string
	switch r.body[1] {
	case 0:
		explain = "Unrecognized Critical Record"
	case 1:
		explain = "Bad Request"
	case 2:
		explain = "Internal Server Error"
	}
	buf.WriteString(fmt.Sprintf("Error Code = %d (%s)", r.body[1], explain))
	return nil
}

func printWarning(r *record, buf *bytes.Buffer) error {
	if r.bodyLen != 2 {
		return fmt.Errorf("unexpected body length in Warning record: %d", r.bodyLen)
	}
	code := (int(r.body[0]) << 8) + int(r.body[1])
	buf.WriteString(fmt.Sprintf("Warning Code = %d (0x%02X%02X)", code, r.body[0], r.body[1]))
	return nil
}

func printAEADAlgorithm(r *record, buf *bytes.Buffer) error {
	if r.bodyLen == 0 {
		buf.WriteString("None of given AEAD algorithm supported")
		return nil
	}
	if r.bodyLen != 2 {
		return fmt.Errorf("unexpected body length in Warning record: %d", r.bodyLen)
	}
	if r.body[0] != 0x00 || r.body[1] == 0x00 || r.body[1] > 0x21 {
		return fmt.Errorf("unrecognized AEAD algorithm ID in AEAD Algorhithm record: %d",
			(int(r.body[0])<<8)+int(r.body[1]))
	}
	buf.WriteString(fmt.Sprintf("Algorithm = %s", datastruct.GetAEADName(r.body[1])))
	return nil
}

func printCookie(r *record, buf *bytes.Buffer) error {
	if r.bodyLen <= 16 {
		return fmt.Errorf("too short body length in Cookie record: %d", r.bodyLen)
	}
	buf.WriteString("Cookie = 0x")
	for _, b := range r.body[:16] {
		buf.WriteString(fmt.Sprintf("%02X", b))
	}
	buf.WriteString("...")
	return nil
}

func printServer(r *record, buf *bytes.Buffer) error {
	if r.bodyLen == 0 {
		return errors.New("empty body in NTPv4 Server record")
	}
	buf.WriteString("Server = ")
	buf.Write(r.body)
	return nil
}

func printPort(r *record, buf *bytes.Buffer) error {
	if r.bodyLen != 2 {
		return fmt.Errorf("unexpected body length in NTPv4 Port record: %d", r.bodyLen)
	}
	port := (int(r.body[0]) << 8) + int(r.body[1])
	buf.WriteString(fmt.Sprintf("Port = %d", port))
	return nil
}
