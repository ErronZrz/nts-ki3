package async

import (
	"active/addr"
	"active/datastruct"
	"active/utils"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net"
	"sync"
	"time"
)

const (
	configPath           = "../resources/"
	localPortKey         = "async.local_port"
	checkIntervalKey     = "async.read.check_interval"
	timeoutKey           = "async.read.timeout"
	haltTimeKey          = "async.send.halt_time"
	partsKey             = "async.send.parts"
	defaultLocalPort     = 11123
	defaultCheckInterval = 1000
	defaultTimeout       = 5000
	defaultHaltTime      = 0
	defaultParts         = 1
)

var (
	checkInterval time.Duration
	timeout       time.Duration
	haltTime      time.Duration
	parts         int
	localPort     int
	errCh         chan error
	dataCh        chan *datastruct.RcvPayload
	wg            *sync.WaitGroup
)

func init() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("properties")
	viper.SetDefault(localPortKey, defaultLocalPort)
	viper.SetDefault(checkIntervalKey, defaultCheckInterval)
	viper.SetDefault(timeoutKey, defaultTimeout)
	viper.SetDefault(haltTimeKey, defaultHaltTime)
	viper.SetDefault(partsKey, defaultParts)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("err reading resource file: %v", err)
		return
	}
	localPort = viper.GetInt(localPortKey)
	checkInterval = time.Duration(viper.GetInt64(checkIntervalKey)) * time.Millisecond
	timeout = time.Duration(viper.GetInt64(timeoutKey)) * time.Millisecond
	haltTime = time.Duration(viper.GetInt64(haltTimeKey)) * time.Millisecond
	parts = viper.GetInt(partsKey)
}

func DialNetworkNTP(cidr string) <-chan *datastruct.RcvPayload {
	errCh = make(chan error)
	finishCh := make(chan struct{})
	go func(finishCh <-chan struct{}, errCh <-chan error) {
		for {
			select {
			case <-finishCh:
				return
			case err := <-errCh:
				fmt.Println(err)
			}
		}
	}(finishCh, errCh)

	wg = new(sync.WaitGroup)
	wg.Add(parts)
	dataCh = make(chan *datastruct.RcvPayload, 1024)

	go func(wg *sync.WaitGroup, finishCh chan<- struct{}) {
		wg.Wait()
		finishCh <- struct{}{}
		close(finishCh)
		close(dataCh)
	}(wg, finishCh)

	networks := utils.SplitCIDR(cidr, parts)
	timeBetweenParts := haltTime / time.Duration(parts)

	for i := 0; i < parts; i++ {
		go singleWriteRead(networks[i], i)
		<-time.After(timeBetweenParts)
	}

	return dataCh
}

func singleWriteRead(cidr string, index int) {
	ctx, cancel := context.WithCancel(context.Background())
	localAddr := &net.UDPAddr{Port: localPort + index}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		errCh <- err
		cancel()
		wg.Done()
		return
	}

	doneCh := make(chan struct{})
	go writeNetWorkNTP(cidr, conn, doneCh)
	go readNetworkNTP(ctx, cidr, conn, doneCh)

	go func() {
		<-doneCh
		time.After(timeout)
		cancel()
		<-doneCh
		close(doneCh)
		wg.Done()
		<-time.After(time.Second)
		_ = conn.Close()
	}()
}

func writeNetWorkNTP(cidr string, conn *net.UDPConn, doneCh chan<- struct{}) {
	defer func() {
		doneCh <- struct{}{}
	}()

	generator, err := addr.NewModuloGenerator(cidr)
	if err != nil {
		errCh <- err
		return
	}
	for generator.HasNext() {
		probeNext(generator.NextHost(), conn)
		if haltTime > 0 {
			<-time.After(haltTime)
		}
	}
}

func probeNext(host string, conn *net.UDPConn) {
	remoteAddr, err := net.ResolveUDPAddr("udp", host+":123")
	if err != nil {
		errCh <- err
		return
	}

	_, err = conn.WriteToUDP(utils.VariableData(), remoteAddr)
	if err != nil {
		errCh <- err
		return
	}
}
