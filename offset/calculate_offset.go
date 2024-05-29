package offset

import (
	"active/parser"
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func CalculateOffsets(path string) error {
	// 创建文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	// 创建结果文件
	resultFile, err := os.Create(path[:len(path)-4] + "_offset.txt")
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
		_ = resultFile.Close()
	}()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(resultFile)
	errCh := make(chan error, 16)

	go func() {
		for err := range errCh {
			fmt.Println(err)
		}
	}()

	var num, finished int

	for scanner.Scan() {
		num++
		fmt.Printf("%d (%d)\n", num, finished)
		line := scanner.Text()
		ip := strings.Split(line, "\t")[0]

		parser.CookieMap = make(map[byte][]byte)
		parser.TheHost = ""
		parser.ThePort = "123"
		ServerTimestampsMap[0x00] = new(ServerTimestamps)
		ServerTimestampsMap[0x0F] = new(ServerTimestamps)
		ServerTimestampsMap[0x10] = new(ServerTimestamps)
		ServerTimestampsMap[0x11] = new(ServerTimestamps)

		wg := new(sync.WaitGroup)
		wg.Add(3)
		go RecordNTSTimestamps(ip, 0x0F, wg, errCh)
		go RecordNTSTimestamps(ip, 0x10, wg, errCh)
		go RecordNTSTimestamps(ip, 0x11, wg, errCh)
		wg.Wait()

		// 如果没有完成 NTS 通信则跳过
		if ServerTimestampsMap[0x0F].T1.IsZero() {
			continue
		}

		if parser.TheHost != "" {
			ip = parser.TheHost
		}
		wg.Add(1)
		go RecordNTPTimestamps(ip+":"+parser.ThePort, wg, errCh)
		wg.Wait()

		_, err = writer.WriteString(generateLine(ip))
		if err != nil {
			return err
		}
		finished++
	}

	return writer.Flush()
}

func generateLine(ip string) string {
	info1 := ServerTimestampsMap[0x00]
	info2 := ServerTimestampsMap[0x0F]
	info3 := ServerTimestampsMap[0x10]
	info4 := ServerTimestampsMap[0x11]
	offset1 := getOffset(info1.T1, info1.T2, info1.T3, info1.T4)
	offset2 := getOffset(info2.T1, info2.T2, info2.T3, info2.T4)
	offset3 := getOffset(info2.RealT1, info2.T2, info2.T3, info2.T4)
	offset4 := getOffset(info3.T1, info3.T2, info3.T3, info3.T4)
	offset5 := getOffset(info3.RealT1, info3.T2, info3.T3, info3.T4)
	offset6 := getOffset(info4.T1, info4.T2, info4.T3, info4.T4)
	offset7 := getOffset(info4.RealT1, info4.T2, info4.T3, info4.T4)
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		ip, offset1, offset2, offset3, offset4, offset5, offset6, offset7,
	)
}

func getOffset(t1, t2, t3, t4 time.Time) string {
	if t1.IsZero() || t2.IsZero() || t3.IsZero() || t4.IsZero() {
		return "-"
	}
	offset := (t2.Sub(t1.UTC()) + t3.Sub(t4.UTC())) / 2
	return fmt.Sprintf("%.3f", float64(offset.Nanoseconds()/1000)/1000)
}
