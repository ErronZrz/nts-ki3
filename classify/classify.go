package classify

import (
	"bytes"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func FetchNTPPackets(filePath string, packets map[string][][]byte) error {
	// 打开 .pcap 文件
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return err
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// 提取源 ip 地址
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}
		ip, _ := ipLayer.(*layers.IPv4)

		srcIP := ip.SrcIP.String() // 获取源 ip 地址的字符串表示

		if _, ok := packets[srcIP]; ok {
			continue
		}

		// 检查是否是 UDP 数据包
		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if udpLayer == nil {
			continue
		}
		udp, _ := udpLayer.(*layers.UDP)

		// 检查端口是否为 NTP 的默认端口 123
		if udp.SrcPort != 123 && udp.DstPort != 123 {
			continue
		}

		// 将 UDP 数据添加到对应 ip 的列表中
		packets[srcIP] = append(packets[srcIP], udp.Payload)
	}

	return nil
}

func GetField(data []byte, name string) string {
	prefix := []byte(name + "=")
	start := bytes.Index(data, prefix)
	if start == -1 {
		return "" // 如果找不到"version="，直接返回空字符串
	}

	// 查找从"version="之后的第一个引号开始的位置
	startQuote := start + len(prefix)
	if startQuote >= len(data) {
		return "" // 检查边界，确保不会越界
	}

	// 从"version="后的位置开始，找到下一个引号
	endQuote := bytes.IndexByte(data[startQuote+1:], '"')
	if endQuote == -1 {
		return "" // 如果找不到闭合的引号，返回空字符串
	}

	// 提取完整的"version="字符串
	value := data[start : startQuote+endQuote+2]
	return string(value)
}

func ClassifyNTPRequest(p *NTPPacket) string {
	// 若 Stratum 非 0，则为 ntpd
	if p.Stratum != 0 {
		return "ntpd"
	}
	// 若 Root Delay 或 Root Dispersion 介于 0 到 1 之间，则为 ntpd
	if p.RootDelay&(1<<16-1) != 0 || p.RootDisp&(1<<16-1) != 0 {
		return "ntpd"
	}
	// 若 Reference ID 非空，则为 ntpd
	if p.RefID != 0 {
		return "ntpd"
	}
	// 若 Origin Timestamp 非 0 或者 Receive Timestamp 非 0，则为 ntpd
	if p.OriginTimestamp != 0 || p.ReceiveTimestamp != 0 {
		return "ntpd"
	}
	// 若 Reference Timestamp 非 0，则为 w32tm
	if p.RefTimestamp != 0 {
		return "w32tm"
	}
	// 若 Poll 为 3 或 Precision 为 -6 或 Root Delay 等于 1，则为 ntpdate
	if p.Poll == 3 || p.Precision == 0xFA || p.RootDelay == (1<<16) {
		return "ntpdate"
	}
	// 若 Precision 不为 32，说明为其他软件
	if p.Precision != 32 {
		return "other"
	}
	// 若 Poll 为 0，则为 NTPsec，否则为 Chrony
	if p.Poll == 0 {
		return "ntpsec"
	}
	return "chrony"
}
