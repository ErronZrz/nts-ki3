package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	preciseFormat = "2006-01-02 15:04:05.000000 UTC"
)

var (
	startingPoint = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	fixedData     []byte
	variableData  []byte
	secData       []byte
	globalStart   = time.Now()
)

func init() {
	fixedData = []byte{
		0xDB, 0x00, 0x04, 0xFA, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
		// 0x01, 0x02, 0x00, 0x04,
	}
	variableData = []byte{
		0xDB, 0x00, 0x04, 0xFA, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
	}
	secData = []byte{
		0x23, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
	}
}

func FixedData() []byte {
	return fixedData
}

func FromInt8(i int8) string {
	val := math.Pow(2, float64(i))
	scientific := FormatScientific(val)
	return fmt.Sprintf("2^%d (%s) sec", i, scientific)
}

func FormatScientific(f float64) string {
	if f == 0 {
		return "0"
	}
	if f >= 0.001 && f <= 1000 {
		return strconv.FormatFloat(f, 'f', 3, 64)
	}
	exp := int(math.Floor(math.Log10(f)))
	mantissa := f / math.Pow10(exp)
	return fmt.Sprintf("%.3fe%d", mantissa, exp)
}

func CalculateDelay(timestamp []byte, another time.Time) time.Duration {
	t := binary.BigEndian.Uint64(timestamp)
	seconds := int64(t >> 32)
	nano := (int64(t&0xFFFF_FFFF) * int64(time.Second)) >> 32
	d := startingPoint.Add(time.Duration(seconds) * time.Second).Add(time.Duration(nano))
	delay := d.Sub(another)
	return delay
}

func FormatTimestamp(timestamp []byte) string {
	return ParseTimestamp(timestamp).Format(preciseFormat)
}

func ParseTimestamp(timestamp []byte) time.Time {
	intPart := binary.BigEndian.Uint32(timestamp[:4])
	fracPart := binary.BigEndian.Uint32(timestamp[4:])
	intTime := startingPoint.Add(time.Duration(intPart) * time.Second)
	fracDuration := (time.Duration(fracPart) * time.Second) >> 32
	return intTime.Add(fracDuration)
}

func GetTimestamp(t time.Time) []byte {
	d := t.Sub(startingPoint)
	seconds := d / time.Second
	high32 := seconds << 32
	nano := d - seconds*time.Second
	low32 := (nano << 32) / time.Second
	res := make([]byte, 8)
	binary.BigEndian.PutUint64(res, uint64(high32|low32))
	return res
}

func TimestampValue(timestamp []byte) float64 {
	intPart := binary.BigEndian.Uint32(timestamp[:4])
	fracPart := binary.BigEndian.Uint32(timestamp[4:])

	return float64(intPart) + float64(fracPart)/float64(1<<32)
}

func RootDelayToValue(data []byte) float64 {
	val := binary.BigEndian.Uint32(data)
	return float64(val) / (1 << 16)
}

func VariableData() []byte {
	timestamp := GetTimestamp(GlobalNowTime())
	copy(variableData[40:], timestamp)
	return variableData
}

func SecData() []byte {
	timestamp := GetTimestamp(GlobalNowTime())
	copy(secData[40:], timestamp)
	return secData
}

func DurationToStr(t1, t2 time.Time) string {
	d := t2.Sub(t1)
	if d < 0 {
		return "-" + DurationToStr(t2, t1)
	}

	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%dμs", d.Microseconds())
	}
	if d < 10*time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	d = d.Round(time.Second)
	h, m, s := d/time.Hour, (d/time.Minute)%60, (d/time.Second)%60
	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func PrintBytes(data []byte, rowLen int) string {
	buf := new(bytes.Buffer)
	rows := len(data) / rowLen
	for i := 0; i < rows; i++ {
		for _, b := range data[i*rowLen : (i+1)*rowLen] {
			buf.WriteString(fmt.Sprintf("%02X ", b))
		}
		buf.WriteByte('\n')
	}
	if len(data) > rows*rowLen {
		for _, b := range data[rows*rowLen:] {
			buf.WriteString(fmt.Sprintf("%02X ", b))
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func SplitCIDR(cidr string, parts int) []string {
	pow, err := CidrPow(cidr)
	if err != nil {
		fmt.Printf("parse CIDR error: %v\n", err)
	}

	num := 1
	for num > 0 && num < parts {
		num <<= 1
		pow--
	}
	if num != parts {
		fmt.Printf("bad parameter: %d", parts)
	}

	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Printf("parse CIDR error: %v\n", err)
	}
	b := ip.Mask(ipNet.Mask)
	basic := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])

	res := make([]string, parts)
	addrNum := 1 << pow

	for i := 0; i < parts; i++ {
		v := basic + uint32(addrNum*i)
		next := fmt.Sprintf("%d.%d.%d.%d/%d", byte(v>>24), byte(v>>16), byte(v>>8), byte(v), 32-pow)
		res[i] = next
	}

	return res
}

func TranslateCountry(countries []string) []string {
	type transReq struct {
		Tp  string   `json:"trans_type"`
		Src []string `json:"source"`
		Id  string   `json:"request_id"`
		D   bool     `json:"detect"`
	}
	type transRes struct {
		Target []string `json:"target"`
	}

	client := &http.Client{}
	tr := transReq{
		Tp:  "zh2en",
		Src: countries,
		Id:  "erron",
		D:   true,
	}
	data, err := json.Marshal(tr)
	// fmt.Println(string(data))
	if err != nil {
		fmt.Printf("marshal error: %v\n", err)
		return nil
	}

	url := "https://api.interpreter.caiyunai.com/v1/translator"
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		fmt.Printf("request error: %v\n", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Authorization", "token wa9ibv5xpgax8d4ds5gk")

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("request error: %v\n", err)
		return nil
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		fmt.Printf("status code error: %s", res.Status)
		return nil
	}

	text, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("read body error: %v\n", err)
		return nil
	}

	var trRes transRes
	err = json.Unmarshal(text, &trRes)
	if err != nil {
		fmt.Printf("unmarshal error: %v\n", err)
		return nil
	}

	return trRes.Target
}

func SameFourBytes(s1, s2 []byte) string {
	// Ensure that both slices are at least 4 bytes long
	if len(s1) < 4 || len(s2) < 4 {
		return ""
	}

	// Create a map to store 4-byte sequences from s1
	sequences := make(map[string]string)
	for i := 0; i <= len(s1)-4; i++ {
		seq := s1[i : i+4]
		hexSeq := fmt.Sprintf("%02X", seq)
		sequences[hexSeq] = hexSeq
	}

	// Check for any sequence in s2 that exists in the map
	for i := 0; i <= len(s2)-4; i++ {
		seq := s2[i : i+4]
		hexSeq := fmt.Sprintf("%02X", seq)
		if val, found := sequences[hexSeq]; found {
			return val
		}
	}

	// If no sequence is found, return an empty string
	return ""
}

func GlobalNowTime() time.Time {
	return globalStart.Add(time.Since(globalStart))
}
