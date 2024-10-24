package classify

import (
	"bufio"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"os"
)

func FetchNTPPackets1(pcapPath string) ([][]byte, error) {
	var res [][]byte
	handle, err := pcap.OpenOffline(pcapPath)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// 检查 UDP
		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if udpLayer == nil {
			continue
		}

		udp, _ := udpLayer.(*layers.UDP)

		// 检查是否至少有一个 123 端口
		if udp.SrcPort != 123 && udp.DstPort != 123 {
			continue
		}

		// 添加结果
		res = append(res, udp.Payload)
	}

	return res, nil
}

func ExtractMode6(data [][]byte, outPath string) error {
	// 打开文件
	file, err := os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	writer := bufio.NewWriter(file)

	for _, d := range data {
		// 检查模式是否为 6，错误位是否为 0
		if len(d) <= 12 || d[0]&0x07 != 6 || d[1]&0x40 != 0 {
			continue
		}
		_, err = writer.Write(d[12:])
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
