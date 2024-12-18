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
	avgOffset := 0.25
	avgDelay := 0.4
	avgDispersion := 0.1
	avgJitter := 0.04
	avgBasicDistance := 0.6
	for i := range data {
		p := &Peer{
			IP:           "192.168.0." + strconv.Itoa(i+4),
			Offset:       avgOffset * (1 + 0.1*rand.NormFloat64()),
			Delay:        avgDelay * (1 + 0.1*rand.NormFloat64()),
			Dispersion:   avgDispersion * (1 + 0.1*rand.NormFloat64()),
			Jitter:       avgJitter * (1 + 0.1*rand.NormFloat64()),
			RootDistance: avgBasicDistance * (1 + 0.1*rand.NormFloat64()),
		}
		p.RootDistance += p.Delay/2 + p.Dispersion + p.Jitter
		data[i] = p
	}
}

func TestClusterAlgorithm(t *testing.T) {
	maxSurvivors := 10
	result := ClusterAlgorithm(data, maxSurvivors)
	for _, p := range result {
		fmt.Printf("%+v\n", p)
	}
}
