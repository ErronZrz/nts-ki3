package clock

import (
	"active/congrat1"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	path := filepath.Join(congrat1.BaseDir, "samples", fmt.Sprintf("%d.txt", congrat1.CurrentBatchID))
	err := writePeers(path, peers)
	if err != nil {
		fmt.Printf("error writing to file: %v", err)
	}
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

// SelectSamples 使用三点排序法筛选至少 minCandidates 数量的 truechimers
// 其中 strict 为 true 时，严格要求样本偏差值在 [low, high] 区间内
// 当 strict 为 false 时，只需偏差值的浮动范围与 [low, high] 区间有交集即可
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
	var (
		count, crossMid int
		low, high       float64
		meet            bool
	)
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
			fmt.Printf("nFake = %d\n", nFake)
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
		if s.offset+k*s.rootDistance > low && s.offset-k*s.rootDistance < high {
			res = append(res, s)
		}
	}
	if len(res) < minCandidates {
		return nil
	}
	return res
}

func writePeers(path string, peers []*Peer) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	writer := bufio.NewWriter(file)
	for _, peer := range peers {
		_, err = writer.WriteString(fmt.Sprintf("%.10f\t%.10f\n", peer.Offset, peer.RootDistance))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
