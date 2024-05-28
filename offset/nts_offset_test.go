package offset

import (
	"active/parser"
	"fmt"
	"sync"
	"testing"
)

func TestRecordNTSTimestamps(t *testing.T) {
	wg := new(sync.WaitGroup)
	errCh := make(chan error)
	parser.CookieMap = make(map[byte][]byte)
	parser.TheHost = ""
	parser.ThePort = "123"
	wg.Add(3)
	go RecordNTSTimestamps("131.130.251.110", 0x0F, wg, errCh)
	go RecordNTSTimestamps("131.130.251.110", 0x10, wg, errCh)
	go RecordNTSTimestamps("131.130.251.110", 0x11, wg, errCh)
	go func() {
		for err := range errCh {
			t.Error(err)
		}
	}()
	wg.Wait()
	close(errCh)
	format := "2006-01-02 15:04:05.000000000"

	info := ServerTimestampsMap[0x0F]
	fmt.Println(info.T1.Format(format))
	fmt.Println(info.RealT1.UTC().Format(format))
	fmt.Println(info.T2.Format(format))
	fmt.Println(info.T3.Format(format))
	fmt.Println(info.T4.UTC().Format(format))

	offset := (info.T2.Sub(info.T1) + info.T3.Sub(info.T4)) / 2
	realOffset := (info.T2.Sub(info.RealT1) + info.T3.Sub(info.T4)) / 2
	fmt.Println(offset)
	fmt.Println(realOffset)

	info = ServerTimestampsMap[0x10]
	fmt.Println(info.T1.Format(format))
	fmt.Println(info.RealT1.UTC().Format(format))
	fmt.Println(info.T2.Format(format))
	fmt.Println(info.T3.Format(format))
	fmt.Println(info.T4.UTC().Format(format))

	offset = (info.T2.Sub(info.T1) + info.T3.Sub(info.T4)) / 2
	realOffset = (info.T2.Sub(info.RealT1) + info.T3.Sub(info.T4)) / 2
	fmt.Println(offset)
	fmt.Println(realOffset)

	info = ServerTimestampsMap[0x11]
	fmt.Println(info.T1.Format(format))
	fmt.Println(info.RealT1.UTC().Format(format))
	fmt.Println(info.T2.Format(format))
	fmt.Println(info.T3.Format(format))
	fmt.Println(info.T4.UTC().Format(format))

	offset = (info.T2.Sub(info.T1) + info.T3.Sub(info.T4)) / 2
	realOffset = (info.T2.Sub(info.RealT1) + info.T3.Sub(info.T4)) / 2
	fmt.Println(offset)
	fmt.Println(realOffset)
}
