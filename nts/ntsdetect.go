package nts

import (
	"active/datastruct"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

func DetectNTSServer(host, serverName string) (*datastruct.NTSDetectPayload, error) {
	//var mutex sync.Mutex
	config := new(tls.Config)
	config.NextProtos = []string{alpnID}
	if serverName != "" {
		config.ServerName = serverName
	} else {
		config.InsecureSkipVerify = true
	}

	addr := host + ":4460"
	//Timeout := 5 * time.Second
	dialer := &net.Dialer{
		Timeout:   20 * time.Second,
		KeepAlive: 20 * time.Second,
	}
	//timeout := 30 * time.Second
	//conn, err := net.DialTimeout("tcp", addr, timeout)
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, config)
	if err != nil {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
		if err != nil {
			conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
			if err != nil {
				conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
				if err != nil {
					conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
					if err != nil {
						conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
						if err != nil {
							conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
							if err != nil {
								conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
								if err != nil {
									conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
									if err != nil {
										conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
										if err != nil {
											conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
											if err != nil {
												return nil, fmt.Errorf("cannot dial TLS server %s: %v", addr, err)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(20 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(20 * time.Second))

	info := datastruct.DetectInfo{
		AEADList:      make([]bool, 34),
		ServerPortSet: make(map[string]struct{}),
	}

	res := &datastruct.NTSDetectPayload{
		Host:   host,
		Port:   4460,
		Secure: !config.InsecureSkipVerify,
		Info:   info,
	}

	//certs := conn.ConnectionState().PeerCertificates
	//if len(certs) > 0 {
	//	res.CertDomain = certs[0].Subject.CommonName
	//}

	for id := byte(0x0F); id <= 0x11; id++ {
		//var conn *tls.Conn
		if id != 0x0F {
			conn, err = newConnection(addr, config, dialer)
			if err != nil {
				return nil, err
			}
		}

		err = singleReadWrite(id, conn, info)
		if err != nil {
			return nil, err
		}

		//name := datastruct.GetAEADName(id) + ":"
		//status := "x"
		//if info.AEADList[id] {
		//	status = "supported"
		//}
		//fmt.Printf("- (%02X) %-27s   %s\n", id, name, status)
	}

	supportOther, err := checkOtherThanAesSivCmac(conn, info)
	if err != nil {
		return nil, err
	}
	if !supportOther {
		//	fmt.Print("- Other AEAD algorithms are not supported\n\n")
		return res, nil
	}

	for id := byte(0x01); id <= 0x21; id++ {
		if id == 0x0F || id == 0x10 || id == 0x11 {
			continue
		}
		var conn *tls.Conn
		conn, err = newConnection(addr, config, dialer)
		if err != nil {
			return nil, err
		}

		err = singleReadWrite(id, conn, info)
		if err != nil {
			return nil, err
		}

		//name := datastruct.GetAEADName(id) + ":"
		//status := "x"
		//if info.AEADList[id] {
		//	status = "supported"
		//}
		//fmt.Printf("- (%02X) %-27s   %s\n", id, name, status)
	}
	fmt.Println()

	return res, nil
}

func newConnection(addr string, config *tls.Config, dialer *net.Dialer) (*tls.Conn, error) {

	conn, err := tls.DialWithDialer(dialer, "tcp", addr, config)
	if err != nil {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
		if err != nil {
			conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
			if err != nil {
				conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
				if err != nil {
					conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
					if err != nil {
						conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
						if err != nil {
							conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
							if err != nil {
								conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
								if err != nil {
									conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
									if err != nil {
										conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
										if err != nil {
											conn, err = tls.DialWithDialer(dialer, "tcp", addr, config)
											if err != nil {
												return nil, fmt.Errorf("cannot dial TLS server %s: %v", addr, err)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return conn, nil
}

