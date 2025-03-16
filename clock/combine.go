package clock

import (
	"fmt"
	"math"
	"slices"
)

type SystemClock struct {
	Offset, Jitter, RootDelay, RootDispersion, PPrev float64
}

var GlobalSystemClock = new(SystemClock)

func CombineAlgorithm(peers []*Peer, selectionJitter float64, useKalman bool) *SystemClock {
	if len(peers) == 0 {
		panic("NO PEERS FOR COMBINATION!!!")
	}

	slices.SortFunc(peers, func(a, b *Peer) int {
		if a.RootDistance < b.RootDistance {
			return -1
		}
		return 1
	})

	var totalWeight, peerJitter, rttError float64
	var offset float64
	sysPeer := peers[0]
	rootDelay := sysPeer.Delay + sysPeer.RootDelay
	for _, p := range peers {
		w := 1.0 / p.RootDistance
		totalWeight += w
		offset += w * p.Offset
		rttError += w * p.RttError
		fmt.Printf("rttError += (%.10f * %.10f = %.10f)\n", w, p.RttError, w*p.RttError)
		peerJitter += w * math.Pow(p.Offset-sysPeer.Offset, 2)
	}
	offset /= totalWeight
	pNow := 1.0
	if useKalman {
		rttError /= totalWeight
		offset, pNow = KalmanFilter(GlobalSystemClock.Offset, offset, GlobalSystemClock.PPrev, rttError)
	}
	jitter := math.Sqrt(math.Pow(selectionJitter, 2) + peerJitter/totalWeight)
	rootDispersion := sysPeer.RootDispersion + sysPeer.Dispersion
	rootDispersion += math.Sqrt(math.Pow(jitter, 2) + math.Pow(sysPeer.Jitter, 2))
	// 理论上这里还需要加一个 PHI * (t4 - t4')，但是太麻烦了先不加
	// 这里的计算公式确实是 offset 的绝对值，但是模拟实验时 offset 是不变的，所以应该减去旧值计算
	// rootDispersion += math.Abs(offset)
	rootDispersion += math.Abs(offset - GlobalSystemClock.Offset)

	return &SystemClock{
		Offset:         offset,
		Jitter:         jitter,
		RootDelay:      rootDelay,
		RootDispersion: rootDispersion,
		PPrev:          pNow,
	}
}
