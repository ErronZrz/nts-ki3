package clock

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

var data []*Peer

func init() {
	data = make([]*Peer, 100)
	avgOffset := -0.3
	avgDelay := 0.1
	avgRootDelay := 0.2
	avgDispersion := 0.1
	avgRootDispersion := 0.2
	avgJitter := 0.04
	for i := range data {
		p := &Peer{
			IP:             "192.168.0." + strconv.Itoa(i+4),
			Offset:         avgOffset * (1 + 0.1*rand.NormFloat64()),
			Delay:          avgDelay * (1 + 0.1*rand.NormFloat64()),
			RootDelay:      avgRootDelay * (1 + 0.1*rand.NormFloat64()),
			Dispersion:     avgDispersion * (1 + 0.1*rand.NormFloat64()),
			RootDispersion: avgRootDispersion * (1 + 0.1*rand.NormFloat64()),
			Jitter:         avgJitter * (1 + 0.1*rand.NormFloat64()),
		}
		p.RootDistance = (p.Delay+p.RootDelay)/2 + p.Dispersion + p.RootDispersion + p.Jitter
		data[i] = p
	}
}

func TestClusterAlgorithm(t *testing.T) {
	maxSurvivors := 40
	result, sj := ClusterAlgorithm(data, maxSurvivors)
	for _, p := range result {
		fmt.Printf("%+v\n", p)
	}
	fmt.Printf("Selection Jitter: %f\n", sj)
}
