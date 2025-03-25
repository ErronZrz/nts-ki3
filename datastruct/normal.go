package datastruct

import (
	"active/region"
	"active/utils"
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type RcvPayload struct {
	Host     string
	Port     int
	Err      error
	Len      int
	SendTime time.Time
	RcvTime  time.Time
	RcvData  []byte
}

func (p *RcvPayload) Print() {
	if p.Err != nil {
		fmt.Println(p.Err)
	} else {
		fmt.Printf(p.Lines())
	}
}

func (p *RcvPayload) Lines() string {
	s := fmt.Sprintf("%d bytes received from %s:%d (%s):\n", p.Len, p.Host, p.Port, region.GetChineseRegion(p.Host, 3))
	buf := bytes.NewBufferString(s)
	buf.WriteString(utils.PrintBytes(p.RcvData, 16))
	// T2 - T1
	sendDelay := utils.CalculateDelay(p.RcvData[32:40], p.SendTime)
	// T4 - T3
	rcvDelay := -utils.CalculateDelay(p.RcvData[40:48], p.RcvTime)
	avgDelay := (sendDelay + rcvDelay) / 2
	offset := (sendDelay - rcvDelay) / 2
	buf.WriteString(fmt.Sprintf("Send delay:    %s\n", durationToStr(sendDelay)))
	buf.WriteString(fmt.Sprintf("Receive delay: %s\n", durationToStr(rcvDelay)))
	buf.WriteString(fmt.Sprintf("Average delay: %s\n", durationToStr(avgDelay)))
	buf.WriteString(fmt.Sprintf("offset:        %s\n", durationToStr(offset)))
	return buf.String()
}

func durationToStr(d time.Duration) string {
	negative := d < 0
	us := d.Microseconds()
	str := strconv.FormatInt(us, 10)
	n := len(str)
	if n <= 3 || (negative && n <= 4) {
		return str + "μs"
	}
	return str[:n-3] + "." + str[n-3:] + "ms"
}
