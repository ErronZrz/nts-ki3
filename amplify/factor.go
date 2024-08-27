package amplify

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"io"
	"os"
)

func AmplificationFactors(path, outPath string) error {
	// 打开 pcap 文件
	handle, err := pcap.OpenOffline(path)
	if err != nil {
		return err
	}
	defer handle.Close()

	// 创建或打开输出文件
	outFile, err := os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(outFile)

	// 用于统计的字典
	stats := make(map[string]struct {
		packets     int
		totalLength int
	})

	// 包解析器
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// 检查网络层和传输层
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if ipLayer != nil && udpLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			udp, _ := udpLayer.(*layers.UDP)

			// 检查端口和 NTP 模式
			if udp.SrcPort == 123 && len(udp.Payload) > 0 {
				version := udp.Payload[0] & 0x07
				if version < 6 {
					continue
				} else if version == 6 && udp.Payload[1]&0x40 != 0 {
					// 错误位为 1
					continue
				} else if version == 7 && udp.Payload[4]&0xF0 != 0 {
					// 错误字段非 0
					continue
				}
				// 统计信息
				sourceIP := ip.SrcIP.String()
				stat := stats[sourceIP]
				stat.packets++
				// IP 数据报长度
				stat.totalLength += int(ip.Length)
				stats[sourceIP] = stat
			}
		}
	}

	// 写入文件
	for ip, stat := range stats {
		line := fmt.Sprintf("%s\t%d\t%d\n", ip, stat.packets, stat.totalLength)
		if _, err := io.WriteString(outFile, line); err != nil {
			return err
		}
	}

	return nil
}
