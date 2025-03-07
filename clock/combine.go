package clock

import (
	"math"
	"slices"
)

type SystemResult struct {
	Offset, Jitter, RootDelay, RootDispersion float64
}

func CombineAlgorithm(peers []*Peer, selectionJitter float64) *SystemResult {
	slices.SortFunc(peers, func(a, b *Peer) int {
		if a.RootDistance < b.RootDistance {
			return -1
		}
		return 1
	})

	var totalWeight, peerJitter float64
	var offset float64
	sysPeer := peers[0]
	rootDelay := sysPeer.Delay + sysPeer.RootDelay
	for _, p := range peers {
		w := 1.0 / p.RootDistance
		totalWeight += w
		offset += w * p.Offset
		peerJitter += w * math.Pow(p.Offset-sysPeer.Offset, 2)
	}
	offset /= totalWeight
	jitter := math.Sqrt(math.Pow(selectionJitter, 2) + peerJitter/totalWeight)
	rootDispersion := sysPeer.RootDispersion + sysPeer.Dispersion
	rootDispersion += math.Sqrt(math.Pow(jitter, 2) + math.Pow(sysPeer.Jitter, 2))
	// 理论上这里还需要加一个 PHI * (t4 - t4')，但是太麻烦了先不加
	rootDispersion += math.Abs(offset)

	return &SystemResult{
		Offset:         offset,
		Jitter:         jitter,
		RootDelay:      rootDelay,
		RootDispersion: rootDispersion,
	}
}
