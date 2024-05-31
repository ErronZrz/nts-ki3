package nts

import (
	"active/datastruct"
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net"
	"strings"
	"time"
)

const (
	aesSivCmac256   = 0x0F
	alpnID          = "ntske/1"
	exportLabel     = "EXPORTER-network-time-security"
	configPath      = "../resources/"
	timeoutKey      = "nts.dial_timeout"
	haltTimeKey     = "nts.detect.halt_time"
	defaultTimeout  = 5000
	defaultHaltTime = 500
)

var (
	reqBytes = []byte{
		0x80, 0x01, 0x00, 0x02, 0x00, 0x00, 0x80, 0x04, 0x00, 0x02, 0x00, 0x0F, 0x80, 0x00, 0x00, 0x00,
	}
	timeout  time.Duration
	haltTime time.Duration
)

func init() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("properties")
	viper.SetDefault(timeoutKey, defaultTimeout)
	viper.SetDefault(haltTimeKey, defaultHaltTime)
	err := viper.ReadInConfig()
	if err != nil {
		// fmt.Printf("error reading resource file: %v", err)
		return
	}
	timeout = time.Duration(viper.GetInt64(timeoutKey)) * time.Millisecond
	haltTime = time.Duration(viper.GetInt64(haltTimeKey)) * time.Millisecond
}

func DialNTSKE(host, serverName string, aeadID byte) (*datastruct.NTSPayload, error) {
	config := new(tls.Config)
	config.NextProtos = []string{alpnID}
	config.MinVersion = tls.VersionTLS12
	config.CipherSuites = []uint16{tls.TLS_AES_128_GCM_SHA256}
	if serverName != "" {
		config.ServerName = serverName
	} else {
		config.InsecureSkipVerify = true
	}

	dialer := &net.Dialer{Timeout: timeout}

	conn, err := tls.DialWithDialer(dialer, "tcp", host+":4460", config)
	if err != nil {
		return nil, fmt.Errorf("cannot dial TLS server %s: %v", host, err)
	}
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("error closing TLS connection %v\n", err)
		}
	}(conn)

	res := &datastruct.NTSPayload{
		Host:   host,
		Port:   4460,
		Secure: !config.InsecureSkipVerify,
	}

	state := conn.ConnectionState()

	certs := state.PeerCertificates
	if len(certs) > 0 {
		cert := certs[0]
		res.CertDomain = cert.Subject.CommonName
		res.RightIP = checkDNS(res.CertDomain, host)
		res.Expired = time.Now().After(cert.NotAfter)
	}

	if aeadID < 0x01 || aeadID > 0x21 {
		aeadID = aesSivCmac256
	}
	reqBytes[11] = aeadID

	_, err = conn.Write(reqBytes)
	if err != nil {
		return nil, fmt.Errorf("send NTS-KE request failed: %v", err)
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("read NTS-KE response failed: %v", err)
	}

	ctx := make([]byte, 4)
	ctx[3] = aeadID
	keyLength := 32
	if aeadID == 0x10 {
		keyLength = 48
	} else if aeadID == 0x11 {
		keyLength = 64
	}
	res.C2SKey, err = state.ExportKeyingMaterial(exportLabel, append(ctx, 0x00), keyLength)
	if err != nil {
		return nil, fmt.Errorf("export C2S key failed: %v", err)
	}
	res.S2CKey, err = state.ExportKeyingMaterial(exportLabel, append(ctx, 0x01), keyLength)
	if err != nil {
		return nil, fmt.Errorf("export S2C key failed: %v", err)
	}

	res.Len = len(data)
	res.RcvData = data
	return res, nil
}

func checkDNS(domain, ip string) bool {
	if strings.Contains(domain, "*") {
		return true
	}

	ips, err := net.LookupIP(domain)
	if err != nil {
		return false
	}

	for _, resolvedIP := range ips {
		if resolvedIP.String() == ip {
			return true
		}
	}

	return false
}
