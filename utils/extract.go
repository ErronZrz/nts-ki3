package utils

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"os"
)

func ExtractNTPPackets(dataPath, dstDir string) error {
	handle, err := pcap.OpenOffline(dataPath)
	if err != nil {
		return fmt.Errorf("error opening pcap file: %w", err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if udpLayer == nil {
			continue
		}

		udp, _ := udpLayer.(*layers.UDP)
		// NTP typically uses port 123
		port := udp.DstPort
		if port != 123 && port != 4123 && port != 8123 {
			continue
		}

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}

		ip, _ := ipLayer.(*layers.IPv4)

		// Extracting NTP payload
		ntpPayload := udp.Payload

		// Create a filename from the destination IP
		fileName := fmt.Sprintf("%s/%s_%d.pkt", dstDir, ip.DstIP, port)
		err = os.WriteFile(fileName, ntpPayload, 0644)
		if err != nil {
			return fmt.Errorf("error writing NTP data to file: %w", err)
		}
	}

	return nil
}
