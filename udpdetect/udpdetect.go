package udpdetect

import (
	"active/addr"
	"active/datastruct"
	"active/utils"
	"github.com/spf13/viper"
	"net"
	"sync"
	"time"
)

const (
	configPath       = "../resources/"
	timeoutKey       = "detection.rcv_header.timeout"
	batchSizeKey     = "detection.send_udp.batch_size"
	defaultTimeout   = 500
	defaultBatchSize = 256
)

var (
	timeout time.Duration
)

func init() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("properties")
	viper.SetDefault(timeoutKey, defaultTimeout)
	viper.SetDefault(batchSizeKey, defaultBatchSize)
	err := viper.ReadInConfig()
	if err != nil {
		// fmt.Printf("error reading resource file: %v", err)
		return
	}
	milli := time.Duration(viper.GetInt64(timeoutKey))
	if milli == 0 {
		milli = defaultTimeout
	}
	timeout = time.Millisecond * milli
}

func DialNetworkNTPWithBatchSize(cidr string, batchSize int) <-chan *datastruct.RcvPayload {
	generator, err := addr.NewModuloGenerator(cidr)
	if err != nil {
		return nil
	}
	num := generator.TotalNum()
	chSize := 1024
	if num < chSize {
		chSize = num
	}
	dataCh := make(chan *datastruct.RcvPayload, chSize)
	wg := new(sync.WaitGroup)
	// fmt.Printf("Num of addresses: %d\n", num)
	wg.Add(num)
	batchNum := num / batchSize
	for i := 0; i < batchNum; i++ {
		for j := 0; j < batchSize; j++ {
			hostStr := generator.NextHost()
			go writeToAddr(hostStr+":123", dataCh, wg)
		}
		time.Sleep(timeout)
	}
	for generator.HasNext() {
		hostStr := generator.NextHost()
		go writeToAddr(hostStr+":123", dataCh, wg)
	}
	go func() {
		wg.Wait()
		close(dataCh)
	}()
	return dataCh
}

func DialNetworkNTP(cidr string) <-chan *datastruct.RcvPayload {
	return DialNetworkNTPWithBatchSize(cidr, viper.GetInt(batchSizeKey))
}

func writeToAddr(addr string, ch chan<- *datastruct.RcvPayload, wg *sync.WaitGroup) {
	defer wg.Done()
	payload := &datastruct.RcvPayload{Host: addr[:len(addr)-4], Port: 123}
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		payload.Err = err
		ch <- payload
		return
	}
	// fmt.Println(udpAddr.Print())
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		payload.Err = err
		ch <- payload
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	payload.SendTime = time.Now()
	_, err = conn.Write(utils.FixedData())
	if err != nil {
		payload.Err = err
		ch <- payload
		return
	}
	buf := make([]byte, 128)
	err = conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		payload.Err = err
		ch <- payload
		return
	}
	n, _, err := conn.ReadFromUDP(buf)
	if err == nil && n > 0 {
		payload.RcvTime = time.Now()
		payload.Len = n
		payload.RcvData = buf[:n]
		ch <- payload
	}
}
