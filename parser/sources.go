package parser

import "fmt"

var (
	sourceMap map[string]string
)

func init() {
	sourceMap = map[string]string{
		"GOES": "Geosynchronous Orbit Environment Satellite",
		"GPS":  "Global Position System",
		"GAL":  "Galileo Positioning System",
		"PPS":  "Generic pulse-per-second",
		"IRIG": "Inter-Range Instrumentation Group",
		"WWVB": "LF Radio WWVB Ft. Collins, CO 60 kHz",
		"DCF":  "LF Radio DCF77 Mainflingen, DE 77.5 kHz",
		"HBG":  "LF Radio HBG Prangins, HB 75 kHz",
		"MSF":  "LF Radio MSF Anthorn, UK 60 kHz",
		"JJY":  "LF Radio JJY Fukushima, JP 40 kHz, Saga, JP 60 kHz",
		"LORC": "MF Radio LORAN C station, 100 kHz",
		"TDF":  "MF Radio Allouis, FR 162 kHz",
		"CHU":  "HF Radio CHU Ottawa, Ontario",
		"WWV":  "HF Radio WWV Ft. Collins, CO",
		"WWVH": "HF Radio WWVH Kauai, HI",
		"NIST": "NIST telephone modem",
		"ACTS": "NIST telephone modem",
		"USNO": "USNO telephone modem",
		"PTB":  "European telephone modem",
		"PTP":  "Precise Time Protocol",
	}
}

func completeSource(s []byte) string {
	var str string
	if s[3] != 0x00 {
		str = string(s)
	} else {
		str = string(s[:3])
	}
	if complete, ok := sourceMap[str]; ok {
		return fmt.Sprintf("%s (%s)", str, complete)
	}
	return str
}
