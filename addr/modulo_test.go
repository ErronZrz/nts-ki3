package addr

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestModuloGenerator(t *testing.T) {
	for pow := 32; pow >= 8; pow-- {
		cidr := randomCIDR(pow)
		fmt.Printf("\n\nCIDR: %s\n", cidr)
		g, err := NewModuloGenerator(cidr)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("total=%d, root=%d, seed=%d, basic=%d\n", g.total, g.root, g.seed, g.basic)
		num := count(g)
		if num != 1<<(32-pow) {
			t.Errorf("want %d addresses but got %d", 1<<(32-pow), num)
		}
	}
}

func randomCIDR(pow int) string {
	a, b, c, d := rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256)
	return fmt.Sprintf("%d.%d.%d.%d/%d", a, b, c, d, pow)
}

func count(g *ModuloGenerator) int {
	res := 0
	for g.HasNext() {
		res++
		_ = g.NextHost()
	}
	fmt.Printf("used=%d, total=%d, next=%d, seed=%d\n", g.used, g.total, g.next, g.seed)
	return res
}
