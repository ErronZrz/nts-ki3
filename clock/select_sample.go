package clock

import (
	"fmt"
	"sort"
)

type Sample struct {
	ip           string
	offset       float64
	rootDistance float64
	flag         int
}

func (s *Sample) realValue() float64 {
	return s.offset + float64(s.flag)*s.rootDistance
}

func NewSample(peer *Peer) *Sample {
	return &Sample{
		ip:           peer.IP,
		offset:       peer.Offset,
		rootDistance: peer.RootDistance,
	}
}

func SelectPeers(peers []*Peer, minCandidates int, strict bool) []*Peer {
	n := len(peers)
	ip2peer := make(map[string]*Peer)
	samples := make([]*Sample, n)
	for i, peer := range peers {
		ip2peer[peer.IP] = peer
		samples[i] = NewSample(peer)
	}
	result := SelectSamples(samples, minCandidates, strict)
	if result == nil {
		return nil
	}
	res := make([]*Peer, 0)
	for _, s := range result {
		res = append(res, ip2peer[s.ip])
	}

	return res
}

func SelectSamples(samples []*Sample, minCandidates int, strict bool) []*Sample {
	n := len(samples)
	if n < minCandidates {
		return nil
	}
	var extended []*Sample
	for _, s := range samples {
		lower := &Sample{
			ip:           s.ip,
			offset:       s.offset,
			rootDistance: s.rootDistance,
			flag:         -1,
		}
		upper := &Sample{
			ip:           s.ip,
			offset:       s.offset,
			rootDistance: s.rootDistance,
			flag:         1,
		}
		extended = append(extended, lower, s, upper)
	}
	sort.Slice(extended, func(i, j int) bool {
		return extended[i].realValue() < extended[j].realValue()
	})
	var count, crossMid int
	var low, high float64
	var meet bool
	for nFake := 0; nFake < n/2; nFake++ {
		count, crossMid = 0, 0
		for _, s := range extended {
			low = s.realValue()
			count -= s.flag
			if s.flag == 0 {
				crossMid++
			}
			if count >= n-nFake {
				break
			}
		}
		count = 0
		for i := len(extended) - 1; i >= 0; i-- {
			s := extended[i]
			high = s.realValue()
			count += s.flag
			if s.flag == 0 {
				crossMid++
			}
			if count >= n-nFake {
				break
			}
		}
		if (!strict || crossMid <= nFake) && low < high {
			fmt.Printf("nFake=%d ", nFake)
			meet = true
			break
		}
	}
	if !meet {
		return nil
	}
	var res []*Sample
	k := 1.0
	if strict {
		k = 0.0
	}
	for _, s := range samples {
		if s.offset+k*s.rootDistance >= low && s.offset-k*s.rootDistance <= high {
			res = append(res, s)
		}
	}
	if len(res) < minCandidates {
		return nil
	}
	return res
}
