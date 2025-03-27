package clock

import (
	"fmt"
	"math"
	"slices"
)

type SystemClock struct {
	Offset, Cumsum, Jitter, RootDelay, RootDispersion, Skew float64
	PPrev                                                   [2][2]float64
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
	pNow := InitialP
	skew := InitialSkew
	if useKalman {
		rttError /= totalWeight
		prev := KalmanState{GlobalSystemClock.Offset, 0.001, GlobalSystemClock.PPrev}
		next := KalmanFilterSkew(prev, offset, rttError, 600)
		offset = next.Offset
		skew = next.Skew
		pNow = next.P
	}
	jitter := math.Sqrt(math.Pow(selectionJitter, 2) + peerJitter/totalWeight)
	rootDispersion := sysPeer.RootDispersion + sysPeer.Dispersion
	rootDispersion += math.Sqrt(math.Pow(jitter, 2) + math.Pow(sysPeer.Jitter, 2))
	// 理论上这里还需要加一个 PHI * (t4 - t4')，但是太麻烦了先不加
	// 这里的计算公式确实是 offset 的绝对值，但是模拟实验时 offset 是不变的，所以应该减去旧值计算
	// rootDispersion += math.Abs(offset)
	// rootDispersion += math.Abs(offset - GlobalSystemClock.Offset)
	// 然而我发现应该直接在 NewPeer 函数里面就减去上一轮的 offset，所以这里还是直接取绝对值即可
	rootDispersion += math.Abs(offset)
	if useKalman {
		fmt.Printf(" %.6f\n", rootDispersion)
	}

	return &SystemClock{
		Offset:         offset,
		Cumsum:         GlobalSystemClock.Cumsum + offset,
		Jitter:         jitter,
		RootDelay:      rootDelay,
		RootDispersion: rootDispersion,
		Skew:           skew,
		PPrev:          pNow,
	}
}
