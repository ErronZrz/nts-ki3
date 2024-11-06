package sysprocess

import (
	"fmt"
	"sort"
)

type sample struct {
	id    int
	value float64
	width float64
	flag  int
}

func (s *sample) realValue() float64 {
	return s.value + float64(s.flag)*s.width
}

func SelectSamples(samples []*sample, minCandidates int, strict bool) []*sample {
	n := len(samples)
	if n < minCandidates {
		return nil
	}
	var extended []*sample
	for _, s := range samples {
		lower := &sample{
			id:    s.id,
			value: s.value,
			width: s.width,
			flag:  -1,
		}
		upper := &sample{
			id:    s.id,
			value: s.value,
			width: s.width,
			flag:  1,
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
	var res []*sample
	k := 1.0
	if strict {
		k = 0.0
	}
	for _, s := range samples {
		if s.value+k*s.width >= low && s.value-k*s.width <= high {
			res = append(res, s)
		}
	}
	if len(res) < minCandidates {
		return nil
	}
	return res
}
