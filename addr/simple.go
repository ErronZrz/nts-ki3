package addr

import (
	"active/utils"
	"net"
)

type SimpleGenerator struct {
	nextIP net.IP
	ipNet  *net.IPNet
	total  int
	used   int
}

func NewAddrGenerator(cidr string) (*SimpleGenerator, error) {
	pow, err := utils.CidrPow(cidr)
	if err != nil {
		return nil, err
	}
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	host := ip.Mask(ipNet.Mask)
	g := &SimpleGenerator{
		nextIP: host,
		ipNet:  ipNet,
		total:  1 << pow,
		used:   0,
	}
	return g, nil
}

func (g *SimpleGenerator) TotalNum() int {
	return g.total
}

func (g *SimpleGenerator) HasNext() bool {
	return g.used < g.total
}

func (g *SimpleGenerator) NextHost() string {
	if g.used >= g.total {
		return ""
	}
	g.used++
	res := g.nextIP.String()
	inc(g.nextIP)
	return res
}

func inc(ip []byte) {
	for i := 3; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}
