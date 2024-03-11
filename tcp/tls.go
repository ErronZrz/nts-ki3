package tcp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	alpnID  = "ntske/1"
	timeout = 500000000
)

func IsTLSEnabled(host string, port int, serverName string) bool {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	config := new(tls.Config)
	if serverName != "" {
		config.ServerName = serverName
	} else {
		config.InsecureSkipVerify = true
	}
	tlsConn := tls.Client(conn, config)
	err = tlsConn.Handshake()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func(tlsConn *tls.Conn) {
		err := tlsConn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(tlsConn)
	return true
}

func WriteReadTLS(host string, port int, serverName string, req []byte) ([]byte, error) {
	config := new(tls.Config)
	config.NextProtos = []string{alpnID}
	if serverName != "" {
		config.ServerName = serverName
	} else {
		config.InsecureSkipVerify = true
	}

	hostPort := net.JoinHostPort(host, strconv.Itoa(port))
	dialer := &net.Dialer{Timeout: timeout}

	conn, err := tls.DialWithDialer(dialer, "tcp", hostPort, config)
	if err != nil {
		return nil, fmt.Errorf("cannot dial TLS server %s: %v", hostPort, err)
	}
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	_, err = conn.Write(req)
	if err != nil {
		return nil, fmt.Errorf("write to TLS failed: %v", err)
	}

	res, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("read from TLS failed: %v", err)
	}

	return res, nil
}
