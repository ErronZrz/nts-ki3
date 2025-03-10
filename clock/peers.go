package clock

import (
	"math"
	"slices"
)

const (
	PHI = 15e-6
)

type OriginSample struct {
	T1, T2, T3, T4 float64
	Offset         float64
	Delay          float64
	Dispersion     float64
}

type Peer struct {
	IP             string
	Offset         float64
	Delay          float64
	RootDelay      float64
	Dispersion     float64
	RootDispersion float64
	Jitter         float64
	RootDistance   float64
}

func NewOriginSample(t1, t2, t3, t4, p float64) *OriginSample {
	offset := (t2 + t3 - t1 - t4) / 2
	delay := t2 + t4 - t1 - t3
	dispersion := p + PHI*(t4-t1)
	return &OriginSample{
		T1:         t1,
		T2:         t2,
		T3:         t3,
		T4:         t4,
		Offset:     offset,
		Delay:      delay,
		Dispersion: dispersion,
	}
}

func NewPeer(samples []*OriginSample, ip string, rootDelay, rootDispersion, ts float64) *Peer {
	// 删除了之前 len(samples) >= 2 的限制，也就是说一个样本也能算，笑死
	if len(samples) == 0 {
		return nil
	}
	slices.SortFunc(samples, func(s1, s2 *OriginSample) int {
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
		psi += math.Pow(s.Offset-offset0, 2)
	}
	var jitter float64
	if len(samples) > 1 {
		jitter = math.Sqrt(psi) / float64(len(samples)-1)
	} else {
		jitter = math.Abs(offset0 * (delay0 + rootDelay) / 4)
	}
	return &Peer{
		IP:             ip,
		Offset:         offset0,
		Delay:          delay0,
		RootDelay:      rootDelay,
		Dispersion:     epsilon,
		RootDispersion: rootDispersion,
		Jitter:         jitter,
		RootDistance:   (delay0+rootDelay)/2 + epsilon + rootDispersion + jitter,
	}
}
