package classify

import (
	"active/parser"
	"encoding/binary"
	"errors"
)

type NTPPacket struct {
	Leap              uint8
	Version           uint8
	Mode              uint8
	Stratum           uint8
	Poll              uint8
	Precision         uint8
	RootDelay         uint32
	RootDisp          uint32
	RefID             uint32
	RefTimestamp      uint64
	OriginTimestamp   uint64
	ReceiveTimestamp  uint64
	TransmitTimestamp uint64
}

func ParseNTPPacket(data []byte) (*NTPPacket, error) {
	if len(data) < parser.HeaderLength {
		return nil, errors.New("data is too short")
	}

	p := new(NTPPacket)

	p.Leap = data[0] >> 6
	p.Version = (data[0] >> 3) & 0b00000111
	p.Mode = data[0] & 0b00000111
	p.Stratum = data[1]
	p.Poll = data[2]
	p.Precision = data[3]
	p.RootDelay = binary.BigEndian.Uint32(data[4:8])
	p.RootDisp = binary.BigEndian.Uint32(data[8:12])
	p.RefID = binary.BigEndian.Uint32(data[12:16])
	p.RefTimestamp = binary.BigEndian.Uint64(data[16:24])
	p.OriginTimestamp = binary.BigEndian.Uint64(data[24:32])
	p.ReceiveTimestamp = binary.BigEndian.Uint64(data[32:40])
	p.TransmitTimestamp = binary.BigEndian.Uint64(data[40:48])

	return p, nil
}
