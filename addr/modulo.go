package addr

import (
	"active/utils"
	"fmt"
	"net"
)

type ModuloGenerator struct {
	total int
	used  int
	root  int
	seed  int
	next  int
	basic int
}

func NewModuloGenerator(cidr string) (*ModuloGenerator, error) {
	pow, err := utils.CidrPow(cidr)
	if err != nil {
		return nil, err
	}
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	root, err := utils.SmallestPrime(pow)
	if err != nil {
		return nil, err
	}
	seed, err := utils.GetSeed(pow)
	if err != nil {
		return nil, err
	}
	total := 1 << pow
	host := ip.Mask(ipNet.Mask)
	basic := int(host[0])<<24 | int(host[1])<<16 | int(host[2])<<8 | int(host[3])
	g := &ModuloGenerator{
		total: total,
		used:  0,
		root:  root,
		seed:  seed,
		next:  seed,
		basic: basic,
	}
	return g, nil
}

func (g *ModuloGenerator) TotalNum() int {
	return g.total
}

func (g *ModuloGenerator) HasNext() bool {
	if g.used == 0 {
		return true
	}
	return g.used < g.total && g.next != g.seed
}

func (g *ModuloGenerator) NextHost() string {
	if !g.HasNext() {
		return ""
	}
	if g.used == g.total>>1 {
		g.used++
		return toIPStr(g.basic)
	}
	g.used++
	res := toIPStr(g.basic + g.next)
	for {
		g.next = (g.next * g.seed) % g.root
		if g.next < g.total {
			break
		}
	}
	//if !g.HasNext() {
	//	fmt.Printf("last host: %s\n", res)
	//}
	return res
}

func toIPStr(x int) string {
	return fmt.Sprintf("%d.%d.%d.%d", 0xFF&(x>>24), 0xFF&(x>>16), 0xFF&(x>>8), 0xFF&x)
}
