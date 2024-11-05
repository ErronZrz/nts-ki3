package sysprocess

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestSelectSamples(t *testing.T) {
	valueAvg := 100.0
	widthAvg := 10.0
	// value 服从标准差为 10 的正态分布，width 服从标准差为 2 的正态分布
	samples := make([]*sample, 500)
	for j := 0; j < 10; j++ {
		for i := range samples {
			value := valueAvg + rand.NormFloat64()*10
			width := widthAvg + rand.NormFloat64()*2
			samples[i] = &sample{
				id:    i,
				value: value,
				width: width,
			}
		}
		fmt.Println(len(SelectSamples(samples, 1, false)), len(SelectSamples(samples, 1, true)))
	}
}
