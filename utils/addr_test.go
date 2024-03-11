package utils

import (
	"fmt"
	"testing"
)

func TestFindSeed(t *testing.T) {
	for i := 0; i <= 32; i++ {
		n, _ := SmallestPrime(i)
		m := FindSeed(n)
		fmt.Printf("i=%d, total=%d, n=%d, m=%d\n", i, 1<<i, n, m)
	}
}
