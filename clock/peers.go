package clock

import (
	"math"
	"slices"
)

const (
	PHI = 15e-6
)

type Sample struct {
	T1, T2, T3, T4 float64
	Offset         float64
	Delay          float64
	Dispersion     float64
}

type Peer struct {
	Offset       float64
	Delay        float64
	Dispersion   float64
	Jitter       float64
	RootDistance float64
}

func NewSample(t1, t2, t3, t4, p float64) *Sample {
	offset := (t2 + t3 - t1 - t4) / 2
	delay := t2 + t4 - t1 - t3
	dispersion := p + PHI*(t4-t1)
	return &Sample{
		T1:         t1,
		T2:         t2,
		T3:         t3,
		T4:         t4,
		Offset:     offset,
		Delay:      delay,
		Dispersion: dispersion,
	}
}

func NewPeer(samples []*Sample, ts float64) *Peer {
	if len(samples) < 2 {
		return nil
	}
	slices.SortFunc(samples, func(s1, s2 *Sample) int {
		if s1.Delay < s2.Delay {
			return -1
		}
		return 1
	})
	offset0 := samples[0].Offset
	delay0 := samples[0].Delay
	var epsilon, psi float64
	weight := 1.0
	for _, s := range samples {
		weight /= 2
		epsilon += weight * (s.Dispersion + PHI*(ts-s.T4))
		diff := s.Offset - offset0
		psi += diff * diff
	}
	jitter := math.Sqrt(psi) / float64(len(samples)-1)
	return &Peer{
		Offset:       offset0,
		Delay:        delay0,
		Dispersion:   epsilon,
		Jitter:       jitter,
		RootDistance: delay0/2 + epsilon,
	}
}
