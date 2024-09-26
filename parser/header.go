package parser

import (
	"active/iputil"
	"active/utils"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

type Header struct {
	Leap              string
	Version           string
	Mode              string
	Stratum           string
	Poll              string
	Precision         string
	RootDelay         string
	RootDisp          string
	RefID             string
	RefTimestamp      string
	OriginTimestamp   string
	ReceiveTimestamp  string
	TransmitTimestamp string
}

type stepFunc func([]byte, *Header) error

const (
	HeaderLength = 48
	InitSource   = "INIT"
)

const (
	normalIndicator = iota
	plusSecondIndicator
	minusSecondIndicator
	unsynchronizedIndicator
)

const (
	symmetricActiveMode = iota + 1
	symmetricPassiveMode
	clientMode
	serverMode
	broadcastMode
	controlMessageMode
)

var (
	parseChain = []stepFunc{
		parseLeapIndicator, parseVersion, parseMode, parseStratum, parsePoll,
		parsePrecision, parseRootDelay, parseRootDisp, parseRefID, parseRefTimestamp,
		parseOriginTimestamp, parseReceiveTimestamp, parseTransmitTimestamp,
	}
)

func ParseHeader(data []byte) (*Header, error) {
	if len(data) < HeaderLength {
		//for _, b := range data {
		//	fmt.Printf("%02X", b)
		//}
		return nil, errors.New(fmt.Sprintf("header length %d less than 48", len(data)))
	}
	res := &Header{}
	var err error
	for _, handler := range parseChain {
		err = handler(data, res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (h *Header) Print() {
	fmt.Print(h.Lines())
}

func (h *Header) Lines() string {
	return fmt.Sprintf("Leap:    %s\nVersion: %s\nMode:    %s\nStratum: %s\nPoll:    %s\nPrecision:       %s\nRoot "+
		"Delay:      %s\nRoot Dispersion: %s\nReference ID:    %s\nReference Timestamp: %s\nOrigin Timestamp:    %s\n"+
		"Receive Timestamp:   %s\nTransmit Timestamp:  %s\n\n\n\n", h.Leap, h.Version, h.Mode, h.Stratum, h.Poll, h.Precision,
		h.RootDelay, h.RootDisp, h.RefID, h.RefTimestamp, h.OriginTimestamp, h.ReceiveTimestamp, h.TransmitTimestamp)
}

func parseLeapIndicator(data []byte, h *Header) error {
	li := (data[0] & 0b11000000) >> 6
	switch li {
	case normalIndicator:
		h.Leap = "No warning"
	case plusSecondIndicator:
		h.Leap = "Last minute of the day has 61 seconds"
	case minusSecondIndicator:
		h.Leap = "Last minute of the day has 59 seconds"
	case unsynchronizedIndicator:
		h.Leap = "Clock unsynchronized"
	default:
		return errors.New("leap indicator unknown error")
	}
	return nil
}

func parseVersion(data []byte, h *Header) error {
	version := (data[0] & 0b00111000) >> 3
	if version < 2 || version > 4 {
		return errors.New(fmt.Sprintf("wrong version number: %d", version))
	}
	h.Version = fmt.Sprintf("NTPv%d", version)
	return nil
}

func parseMode(data []byte, h *Header) error {
	mode := data[0] & 0b00000111
	switch mode {
	case symmetricActiveMode:
		h.Mode = "Symmetric active"
	case symmetricPassiveMode:
		h.Mode = "Symmetric passive"
	case clientMode:
		h.Mode = "Client"
	case serverMode:
		h.Mode = "Server"
	case broadcastMode:
		h.Mode = "Broadcast"
	case controlMessageMode:
		h.Mode = "Control message"
	default:
		return errors.New(fmt.Sprintf("wrong mode number: %d", mode))
	}
	return nil
}

func parseStratum(data []byte, h *Header) error {
	stratum := data[1]
	if stratum > 16 {
		return errors.New(fmt.Sprintf("wrong stratum number: %d", stratum))
	}
	if stratum == 0 {
		h.Stratum = "Not specified"
	} else if stratum == 16 {
		h.Stratum = "Clock unsynchronized"
	} else {
		h.Stratum = strconv.Itoa(int(stratum))
	}
	return nil
}

func parsePoll(data []byte, h *Header) error {
	poll := int8(data[2])
	h.Poll = utils.FromInt8(poll)
	return nil
}

func parsePrecision(data []byte, h *Header) error {
	precision := int8(data[3])
	h.Precision = utils.FromInt8(precision)
	return nil
}

func parseRootDelay(data []byte, h *Header) error {
	val := binary.BigEndian.Uint32(data[4:8])
	floatVal := float64(val) / (1 << 16)
	h.RootDelay = utils.FormatScientific(floatVal) + " sec"
	return nil
}

func parseRootDisp(data []byte, h *Header) error {
	val := binary.BigEndian.Uint32(data[8:12])
	floatVal := float64(val) / (1 << 16)
	h.RootDisp = utils.FormatScientific(floatVal) + " sec"
	return nil
}

func parseRefID(data []byte, h *Header) error {
	if h.Stratum == "1" || string(data[12:16]) == InitSource {
		// Special reference identifier
		h.RefID = completeSource(data[12:16])
	} else {
		// Normal IP address
		ipStr := fmt.Sprintf("%d.%d.%d.%d", data[12], data[13], data[14], data[15])
		h.RefID = fmt.Sprintf("%s (%s)", ipStr, iputil.GetChineseRegion(ipStr, 3))
	}
	return nil
}

func parseRefTimestamp(data []byte, h *Header) error {
	h.RefTimestamp = utils.FormatTimestamp(data[16:24])
	return nil
}

func parseOriginTimestamp(data []byte, h *Header) error {
	h.OriginTimestamp = utils.FormatTimestamp(data[24:32])
	return nil
}

func parseReceiveTimestamp(data []byte, h *Header) error {
	h.ReceiveTimestamp = utils.FormatTimestamp(data[32:40])
	return nil
}

func parseTransmitTimestamp(data []byte, h *Header) error {
	h.TransmitTimestamp = utils.FormatTimestamp(data[40:48])
	return nil
}
