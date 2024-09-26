package datastruct

import (
	"active/iputil"
	"active/utils"
	"bufio"
	"encoding/binary"
	"fmt"
)

type Statistic struct {
	Domain         string
	IP             string
	Country        string
	Stratum        int
	Poll           int
	Precision      int
	Delay          int
	Offset         int
	ProcessingTime int
	RefCountry     string
	RootDelay      int
	RootDisp       int
}

func NewStatistic(p *RcvPayload) *Statistic {
	res := new(Statistic)

	res.IP = p.Host
	res.Country = iputil.GetCountry(p.Host)

	data := p.RcvData
	stratum := data[1]

	res.Stratum = int(stratum)
	res.Poll = int(int8(data[2]))
	res.Precision = int(int8(data[3]))
	res.RootDelay = int(binary.BigEndian.Uint32(data[4:8]))
	res.RootDisp = int(binary.BigEndian.Uint32(data[8:12]))

	sendDelay := utils.CalculateDelay(data[32:40], p.SendTime)
	rcvDelay := -utils.CalculateDelay(data[40:48], p.RcvTime)
	avgDelay := (sendDelay + rcvDelay) / 2
	offset := (sendDelay - rcvDelay) / 2
	res.Delay = int(avgDelay.Microseconds())
	res.Offset = int(offset.Microseconds())
	res.ProcessingTime = int(binary.BigEndian.Uint64(data[40:48]) - binary.BigEndian.Uint64(data[32:40]))

	if stratum == 1 {
		if data[15] == 0x00 {
			res.RefCountry = string(data[12:15])
		} else {
			res.RefCountry = string(data[12:16])
		}
	} else {
		ipStr := fmt.Sprintf("%d.%d.%d.%d", data[12], data[13], data[14], data[15])
		res.RefCountry = iputil.GetCountry(ipStr)
	}

	return res
}

func (s *Statistic) WriteToCSV(writer *bufio.Writer) error {
	_, err := writer.WriteString(fmt.Sprintf("%s,%s,%s,%d,%d,%d,%d,%d,%d,%s,%d,%d\n",
		s.Domain, s.IP, s.Country, s.Stratum, s.Poll, s.Precision, s.Delay, s.Offset,
		s.ProcessingTime, s.RefCountry, s.RootDelay, s.RootDisp))
	if err != nil {
		return fmt.Errorf("error writing statistic to CSV: %v", err)
	}
	return nil
}
