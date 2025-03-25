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
	port          int
	delta         int
	timeOffset    int
	availability  int
	normalDist    string
	refID         = []byte{0, 0, 0, 0}
	startingPoint = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "malign",
		Short: "Minimal NTP Server with configurable behavior",
		Run:   startServer,
	}

	rootCmd.Flags().IntVarP(&port, "port", "p", 123, "Port number to listen on")
	rootCmd.Flags().IntVarP(&delta, "delta", "d", 0, "Artificial delay in milliseconds")
	rootCmd.Flags().IntVarP(&timeOffset, "timeOffset", "t", 0, "Time offset in milliseconds")
	rootCmd.Flags().IntVarP(&availability, "availability", "a", 100, "Probability of responding (%)")
	rootCmd.Flags().StringVarP(&normalDist, "normalDist", "n", "0,0", "Additional delay in format avg,std")

	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
}

func startServer(_ *cobra.Command, _ []string) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer func() { _ = conn.Close() }()

	buf := make([]byte, 48)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil || n < 48 {
			continue
		}

		if rand.Intn(100) >= availability {
			continue
		}

		go handleRequest(conn, remoteAddr, buf[:n])
	}
}

func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, request []byte) {
	delay := calculateDelay()

	if delay > 0 {
		time.Sleep(delay)
	}

	now := time.Now().Add(time.Duration(timeOffset) * time.Millisecond)
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
	copy(response[40:48], getTimestamp(time.Now().Add(time.Duration(timeOffset)*time.Millisecond)))

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

func calculateDelay() time.Duration {
	parts := strings.Split(normalDist, ",")
	avg, _ := strconv.ParseFloat(parts[0], 64)
	std, _ := strconv.ParseFloat(parts[1], 64)

	normalDelay := rand.NormFloat64()*std + avg
	totalDelay := float64(delta) + normalDelay

	return time.Duration(totalDelay * float64(time.Millisecond))
}
