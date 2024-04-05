package classify

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func FetchNTPPackets(filePath string, limit int) ([][]byte, error) {
	var packets [][]byte

	// 打开.pcap文件
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
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

		// 此时 udp.Payload 即为 NTP 的数据部分
		packets = append(packets, udp.Payload)
		if limit > 0 && len(packets) >= limit {
			break
		}
	}

	return packets, nil
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
