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
	"sync"
	"time"
)

var (
	ip              string
	ports           string
	deltaStr        string
	deltaAvg        float64
	deltaStdDev     float64
	timeOffset      string
	availability    int
	availabilityStr string
	deltaStrategy   string
	refID           = []byte{0, 0, 0, 0}
	startingPoint   = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	mutex           sync.Mutex
)

type availabilityRule struct {
	startMin int
	delta    int
	minAvail int
}

type deltaStep struct {
	startMin int
	interval int
	incr     float64
}

func main666() {
	var rootCmd = &cobra.Command{
		Use:   "malign",
		Short: "Minimal NTP Server with configurable behavior",
		Run:   startServer,
	}

	rootCmd.Flags().StringVarP(&ip, "ip", "i", "0.0.0.0", "IP address to bind to")
	rootCmd.Flags().StringVarP(&ports, "ports", "p", "123", "Port number or range to listen on")
	rootCmd.Flags().StringVarP(&deltaStr, "delta", "d", "0,0", "Artificial delay in format avg,std (ms)")
	rootCmd.Flags().StringVarP(&timeOffset, "timeOffset", "t", "0,0", "Time offset in format avg,std (ms)")
	rootCmd.Flags().IntVarP(&availability, "availability", "a", 100, "Initial probability of responding (%)")
	rootCmd.Flags().StringVarP(&availabilityStr, "availabilityStrategy", "A", "0,0,10", "Availability strategy in format startMin,delta,minAvail")
	rootCmd.Flags().StringVarP(&deltaStrategy, "deltaStrategy", "D", "0,10,0", "Delta adjustment strategy string e.g. 20,30,-200/30,30,200")

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

	parts := strings.Split(deltaStr, ",")
	if len(parts) != 2 {
		_, _ = fmt.Fprintln(os.Stderr, "Invalid delta argument format")
		os.Exit(1)
	}
	avg, _ := strconv.ParseFloat(parts[0], 64)
	std, _ := strconv.ParseFloat(parts[1], 64)
	deltaAvg = avg
	deltaStdDev = std

	availRule := parseAvailabilityStrategy(availabilityStr)
	deltaSteps := parseDeltaStrategies(deltaStrategy)

	go adjustAvailability(availRule)
	for _, strategy := range deltaSteps {
		go func(s deltaStep) {
			time.Sleep(time.Duration(s.startMin) * time.Minute)
			interval := time.Duration(s.interval) * time.Minute
			for {
				mutex.Lock()
				deltaAvg += s.incr
				mutex.Unlock()
				time.Sleep(interval)
			}
		}(strategy)
	}

	for _, port := range portList {
		addr := net.UDPAddr{
			Port: port,
			IP:   net.ParseIP(ip),
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
	mutex.Lock()
	delay := time.Duration((rand.NormFloat64()*deltaStdDev + deltaAvg) * float64(time.Millisecond))
	mutex.Unlock()
	offset := randomOffset(timeOffset)

	if delay > 0 {
		time.Sleep(delay)
	}

	now := time.Now().Add(offset)
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
	copy(response[40:48], getTimestamp(time.Now().Add(offset)))

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

func randomOffset(arg string) time.Duration {
	parts := strings.Split(arg, ",")
	if len(parts) != 2 {
		return 0
	}

	avg, _ := strconv.ParseFloat(parts[0], 64)
	std, _ := strconv.ParseFloat(parts[1], 64)

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

func parseAvailabilityStrategy(s string) availabilityRule {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return availabilityRule{0, 0, 10}
	}
	start, _ := strconv.Atoi(parts[0])
	delta, _ := strconv.Atoi(parts[1])
	minAvail, _ := strconv.Atoi(parts[2])
	return availabilityRule{startMin: start, delta: delta, minAvail: minAvail}
}

func adjustAvailability(rule availabilityRule) {
	time.Sleep(time.Duration(rule.startMin) * time.Minute)

	for {
		availability += rule.delta
		if availability > 100 {
			availability = 100
		}
		if availability < rule.minAvail {
			availability = rule.minAvail
		}
		time.Sleep(time.Minute)
	}
}

func parseDeltaStrategies(s string) []deltaStep {
	var result []deltaStep
	strategies := strings.Split(s, "/")
	for _, item := range strategies {
		parts := strings.Split(item, ",")
		if len(parts) != 3 {
			continue
		}
		start, _ := strconv.Atoi(parts[0])
		interval, _ := strconv.Atoi(parts[1])
		incr, _ := strconv.ParseFloat(parts[2], 64)
		result = append(result, deltaStep{startMin: start, interval: interval, incr: incr})
	}
	return result
}
