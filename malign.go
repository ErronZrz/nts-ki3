package main

import (
	"encoding/binary"
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ports         string
	delta         string
	timeOffset    string
	availability  int
	refID         = []byte{0, 0, 0, 0}
	startingPoint = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
)

func main666() {
	var rootCmd = &cobra.Command{
		Use:   "malign",
		Short: "Minimal NTP Server with configurable behavior",
		Run:   startServer,
	}

	rootCmd.Flags().StringVarP(&ports, "ports", "p", "123", "Port number or range to listen on (e.g. 123 or 3001-3010)")
	rootCmd.Flags().StringVarP(&delta, "delta", "d", "0,0", "Artificial delay in format avg,std (ms)")
	rootCmd.Flags().StringVarP(&timeOffset, "timeOffset", "t", "0,0", "Time offset in format avg,std (ms)")
	rootCmd.Flags().IntVarP(&availability, "availability", "a", 100, "Probability of responding (%)")

	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
}

func startServer(_ *cobra.Command, _ []string) {
	portList, err := parsePorts(ports)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Invalid ports argument: %v\n", err)
		os.Exit(1)
	}

	for _, port := range portList {
		addr := net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("0.0.0.0"),
		}

		conn, err := net.ListenUDP("udp", &addr)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to listen on port %d: %v\n", port, err)
			continue
		}

		go func(c *net.UDPConn) {
			defer func() { _ = c.Close() }()
			buf := make([]byte, 48)

			for {
				n, remoteAddr, err := c.ReadFromUDP(buf)
				if err != nil || n < 48 {
					continue
				}

				if rand.Intn(100) >= availability {
					continue
				}

				go handleRequest(c, remoteAddr, buf[:n])
			}
		}(conn)
	}
	select {} // prevent main from exiting
}

func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, request []byte) {
	delay := calculateRandomDuration(delta)
	timeOffset := calculateRandomDuration(timeOffset)

	if delay > 0 {
		time.Sleep(delay)
	}

	now := time.Now().Add(timeOffset)
	response := make([]byte, 48)
	// LI/VN/Mode (0 4 4)
	response[0] = 0x24
	// Stratum
	response[1] = 0x03
	// Poll
	response[2] = 0x04
	// Precision (-24)
	response[3] = 0xE8
	// Root Delay (0.15625)
	copy(response[4:8], []byte{0x00, 0x00, 0x04, 0x00})
	// Root Dispersion (0.234375)
	copy(response[8:12], []byte{0x00, 0x00, 0x06, 0x00})
	// Reference ID
	copy(response[12:16], refID)
	// Reference Timestamp
	copy(response[16:24], getTimestamp(now.Add(-600*time.Second)))
	// Origin Timestamp
	copy(response[24:32], request[40:48])
	// Receive Timestamp
	copy(response[32:40], getTimestamp(now))
	// Transmit Timestamp
	copy(response[40:48], getTimestamp(time.Now().Add(timeOffset)))

	if delay < 0 {
		time.Sleep(-delay)
	}

	_, err := conn.WriteToUDP(response, addr)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
}

func getTimestamp(t time.Time) []byte {
	d := t.Sub(startingPoint)
	seconds := d / time.Second
	high32 := seconds << 32
	nano := d - seconds*time.Second
	low32 := (nano << 32) / time.Second
	res := make([]byte, 8)
	binary.BigEndian.PutUint64(res, uint64(high32|low32))
	return res
}

func calculateRandomDuration(arg string) time.Duration {
	parts := strings.Split(arg, ",")
	if len(parts) != 2 {
		return 0
	}

	avg, err1 := strconv.ParseFloat(parts[0], 64)
	std, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return 0
	}

	value := rand.NormFloat64()*std + avg
	return time.Duration(value * float64(time.Millisecond))
}

func parsePorts(portsStr string) ([]int, error) {
	if strings.Contains(portsStr, "-") {
		parts := strings.Split(portsStr, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid port range")
		}
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || start > end || start <= 0 || end > 65535 {
			return nil, fmt.Errorf("invalid port range")
		}
		var result []int
		for i := start; i <= end; i++ {
			result = append(result, i)
		}
		return result, nil
	} else {
		port, err := strconv.Atoi(portsStr)
		if err != nil || port <= 0 || port > 65535 {
			return nil, fmt.Errorf("invalid port")
		}
		return []int{port}, nil
	}
}
