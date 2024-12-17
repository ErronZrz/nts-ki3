package async

import (
	"active/datastruct"
	"active/parser"
	"active/utils"
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

func readNetworkNTP(ctx context.Context, cidr string, conn *net.UDPConn, doneCh chan<- struct{}) {
	defer func() {
		doneCh <- struct{}{}
	}()
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 128)

	for {
		select {
		case <-ctx.Done():
			// fmt.Println("Done!")
			return
		default:
			err := conn.SetReadDeadline(time.Now().Add(checkInterval))
			if err != nil {
				fmt.Println(err)
				continue
			}
			n, udpAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			if !ipNet.Contains(udpAddr.IP) {
				fmt.Println("ip out of range: " + udpAddr.IP.String())
				continue
			}
			payload := &datastruct.RcvPayload{
				Host:    udpAddr.IP.String(),
				Port:    udpAddr.Port,
				Len:     n,
				RcvTime: time.Now(),
				RcvData: buf[:n],
			}
			if n < parser.HeaderLength {
				payload.Err = errors.New(fmt.Sprintf("header length %d less than 48", n))
			} else {
				payload.SendTime = utils.ParseTimestamp(buf[24:32])
			}
			dataCh <- payload
		}
	}
}
