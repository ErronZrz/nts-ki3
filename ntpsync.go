package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	timeout = 500 * time.Millisecond
	reqData = []byte{
		0x23, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
	}
	startingPoint1 = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	addr string
	dst  string

	rootCmd = &cobra.Command{
		Use:   "ntp-sync",
		Short: "NTP Sync 工具，用于定期向指定的 NTP 服务器同步时间并记录偏移量",
		Run:   RunTimeSync,
	}
)

func init() {
	// 定义命令行参数
	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "a", "", "NTP 服务器地址（例如：192.168.1.1:123）")
	rootCmd.PersistentFlags().StringVarP(&dst, "dst", "d", "", "偏移量存储文件路径（例如：/root/records/offsets.txt）")
}

func main678() {
	// 执行命令行程序
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RunTimeSync(_ *cobra.Command, _ []string) {
	// 打开输出文件，如果文件不存在则创建
	file, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("无法打开文件 %s: %v", dst, err)
	}
	writer := bufio.NewWriter(file)

	// 解析目标地址
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("无法解析 NTP 服务器地址 %s: %v", addr, err)
	}

	// 创建 UDP 客户端
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatalf("无法连接到 NTP 服务器 %s: %v", addr, err)
	}
	var offsets []int64
	defer func() {
		// 程序停止时将所有偏移量写入文件
		offsetsLine := fmt.Sprintf("[%s]\n", formatOffsets(offsets))
		_, err = writer.WriteString(offsetsLine)
		if err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}
		err = writer.Flush()
		if err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}
		fmt.Println("所有偏移量已写入文件")

		_ = conn.Close()
		_ = file.Close()
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// 捕获 SIGINT 信号（Ctrl+C）
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// 定期发送 NTP 请求并计算偏移量
	loopCount := 0
	for {
		select {
		case <-ticker.C:
			// 增加循环计数
			loopCount++

			// 发送 NTP 请求
			req := makeNTPRequest()

			// 发送并接收响应，设置 500ms 超时
			resp, err := sendNTPRequest(conn, req)
			if err != nil {
				// 如果没有响应，打印并继续
				fmt.Printf("第 %d 次请求未响应\n", loopCount)
				continue
			}

			// 计算并记录偏移量
			offset := calculateOffset(resp)
			offsets = append(offsets, offset)
		case <-sigChan:
			// 捕获到 SIGINT 信号时，执行清理操作
			fmt.Println("\n接收到 Ctrl+C 信号，程序将停止...")
			return
		}
	}
}

func makeNTPRequest() []byte {
	copy(reqData[40:], getTimestamp1(time.Now()))
	return reqData
}

func calculateOffset(resp []byte) int64 {
	t4 := time.Now()
	t1 := parseTimestamp1(resp[24:32])
	t2 := parseTimestamp1(resp[32:40])
	t3 := parseTimestamp1(resp[40:48])
	offset := (t2.Sub(t1) + t3.Sub(t4)) / 2
	return offset.Nanoseconds()
}

func getTimestamp1(t time.Time) []byte {
	d := t.Sub(startingPoint1)
	seconds := d / time.Second
	high32 := seconds << 32
	nano := d - seconds*time.Second
	low32 := (nano << 32) / time.Second
	res := make([]byte, 8)
	binary.BigEndian.PutUint64(res, uint64(high32|low32))
	return res
}

func parseTimestamp1(timestamp []byte) time.Time {
	intPart := binary.BigEndian.Uint32(timestamp[:4])
	fracPart := binary.BigEndian.Uint32(timestamp[4:])
	intTime := startingPoint1.Add(time.Duration(intPart) * time.Second)
	fracDuration := (time.Duration(fracPart) * time.Second) >> 32
	return intTime.Add(fracDuration)
}

// 格式化偏移量为指定的精度
func formatOffsets(offsets []int64) string {
	formatted := make([]string, len(offsets))
	for i, offset := range offsets {
		formatted[i] = fmt.Sprintf("%d", offset)
	}
	return strings.Join(formatted, ", ")
}

// 发送 NTP 请求并接收响应
func sendNTPRequest(conn *net.UDPConn, req []byte) ([]byte, error) {
	// 设置超时时间
	err := conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, fmt.Errorf("设置超时失败: %v", err)
	}

	// 发送 NTP 请求
	_, err = conn.Write(req)
	if err != nil {
		return nil, fmt.Errorf("发送 NTP 请求失败: %v", err)
	}

	// 读取响应
	resp := make([]byte, 48) // NTP 响应大小为 48 字节
	n, err := conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("接收 NTP 响应失败: %v", err)
	}

	// 如果读取的字节数不是 48 字节，说明响应不完整
	if n != 48 {
		return nil, errors.New("NTP 响应数据不完整")
	}

	// 返回响应数据
	return resp, nil
}
