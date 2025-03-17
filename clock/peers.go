package clock

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
)

const (
	PHI         = 15e-6
	OffsetError = 0.05
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
	RttError       float64
}

func NewOriginSample(t1, t2, t3, t4, p float64) *OriginSample {
	offset := (t2+t3-t1-t4)/2 - GlobalSystemClock.Cumsum
	// 人为制造高斯噪声误差
	offset += OffsetError * rand.NormFloat64()
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
		// 原先的排序方式是按照 Delay 从小到大，但是我觉得这完全不合理，所以改成按照时间从大到小
		if s1.T1 > s2.T1 {
			return -1
		}
		return 1
	})
	offset0 := samples[0].Offset
	delay0 := samples[0].Delay
	var epsilon, psi, rttError float64
	// 计算 RTT 对称性误差
	if len(samples) == 2 {
		rttError = math.Abs(samples[0].Delay-samples[1].Delay) / 2
	}
	weight := 1.0
	for _, s := range samples {
		weight /= 2
		epsilon += weight * (s.Dispersion + PHI*(ts-s.T4))
		// 这里的抖动值有个问题，就是收集了多个样本的过程中 offset 本身会变化，所以这真的能反映吗？
		psi += math.Pow(s.Offset-offset0, 2)
	}
	var jitter float64
	if len(samples) > 1 {
		// （接上一个注释）所以这里暂时除以 4 以抵消这个影响，后面再看看怎么改吧
		jitter = math.Sqrt(psi) / float64(len(samples)-1) / 4
	} else {
		jitter = math.Abs(delay0+rootDelay) / 4
	}
	fmt.Printf("sample num = %d, jitter = %.10f\n", len(samples), jitter)
	return &Peer{
		IP:             ip,
		Offset:         offset0,
		Delay:          delay0,
		RootDelay:      rootDelay,
		Dispersion:     epsilon,
		RootDispersion: rootDispersion,
		Jitter:         jitter,
		// RootDistance 在 RFC 8915 里面一般是 LAMBDA = EPSILON + DELTA / 2
		// 然而搜索 synchronization distance 的最后一个结果又加上了抖动，这就很难搞了
		// 又考虑到目前的抖动计算不太规范，所以就先去掉
		RootDistance: (delay0+rootDelay)/2 + epsilon + rootDispersion,
		RttError:     rttError + 0.3,
	}
}
