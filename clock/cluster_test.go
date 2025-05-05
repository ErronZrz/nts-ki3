package clock

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var data = make([]*Peer, 40)

func TestClusterAlgorithm(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	shuffle(-300, r)
	minSurvivors := 3
	result, sj := ClusterAlgorithm(data, minSurvivors)
	for _, p := range result {
		fmt.Printf("%+v\n", p)
	}
	fmt.Printf("Selection Jitter: %f\n", sj)
}

func TestVisualize(t *testing.T) {
	r := rand.New(rand.NewSource(44))
	avgOffset := 0.0
	for i := 0; i < 10; i++ {
		shuffle(avgOffset, r)
		once(i)
	}
}

func once(now int) {
	m := make(map[string]int, len(data))
	for i, p := range data {
		m[p.IP] = i
	}
	offsets := make([]float64, len(data))
	weights := make([]float64, len(data))
	var sum, weightValue, weightSum, finalVal, finalSum, finalWeightSum float64
	for i, p := range data {
		offsets[i] = p.Offset
		weights[i] = 1 / p.RootDistance
		sum += p.Offset
		weightValue += p.Offset * weights[i]
		weightSum += weights[i]
	}
	median := (offsets[len(offsets)/2-1] + offsets[len(offsets)/2]) / 2
	minSurvivors := 3
	result, _ := ClusterAlgorithm(data, minSurvivors)
	resultIndexes := make([]int, len(result))
	for i, p := range result {
		idx := m[p.IP]
		resultIndexes[i] = idx
		finalSum += p.Offset
		finalVal += p.Offset * weights[idx]
		finalWeightSum += weights[idx]
	}
	finalSum /= float64(len(result))
	kMeansRes := KMeans(offsets, 3)
	sort.Float64s(kMeansRes)
	printFloats(offsets, "offsets", now)
	// printFloats(weights, "weights", now)
	printInts(resultIndexes, "selected", now)
	avg := sum / float64(len(offsets))
	weightedAvg := weightValue / weightSum
	center := finalVal / finalWeightSum
	fmt.Printf("avg[%d] = %.10f\nweighted_avg[%d] = %.10f\nmedian[%d] = %.10f\n"+
		"kmeans[%d] = %.10f\ncenter[%d] = %.10f\n",
		now, avg, now, weightedAvg, now, median, now, kMeansRes[1], now, center)
	avgDiff := math.Abs(avg - finalSum)
	weightedAvgDiff := math.Abs(weightedAvg - finalSum)
	medianDiff := math.Abs(median - finalSum)
	kMeansDiff := math.Abs(kMeansRes[1] - finalSum)
	centerDiff := math.Abs(center - finalSum)
	fmt.Printf("diff_data[%d] = np.array([%.10f, %.10f, %.10f, %.10f, %.10f])\n",
		now, avgDiff, weightedAvgDiff, medianDiff, kMeansDiff, centerDiff)
}

func shuffle(avgOffset float64, rd *rand.Rand) {
	// 以下数值单位均为 ms
	avgDelay := 80.0
	avgRootDelay := 160.0
	avgDispersion := 100.0
	avgRootDispersion := 200.0
	avgJitter := 2.4
	for i := range data {
		alpha := 1.0
		v := rd.Float64()
		if v < 0.1 {
			alpha = 0.3
		} else if v < 0.4 {
			alpha = 0.9
		}
		p := &Peer{
			IP:             "192.168.0." + strconv.Itoa(i+4),
			Offset:         avgOffset + 2*rd.NormFloat64(),
			Delay:          avgDelay * (1 + 0.1*rd.NormFloat64()) * alpha,
			RootDelay:      avgRootDelay * (1 + 0.1*rd.NormFloat64()) * alpha,
			Dispersion:     avgDispersion * (1 + 0.1*rd.NormFloat64()) * alpha,
			RootDispersion: avgRootDispersion * (1 + 0.1*rd.NormFloat64()) * alpha,
			Jitter:         avgJitter * (1 + 0.1*rd.NormFloat64()) * alpha,
		}
		p.RootDistance = (p.Delay+p.RootDelay)/2 + p.Dispersion + p.RootDispersion + p.Jitter
		data[i] = p
	}
}

func printFloats(a []float64, name string, idx int) {
	strList := make([]string, len(a))
	for i, v := range a {
		strList[i] = strconv.FormatFloat(v, 'f', 10, 64)
	}
	fmt.Printf("%s[%d] = np.array([%s])\n", name, idx, strings.Join(strList, ", "))
}

func printInts(a []int, name string, idx int) {
	strList := make([]string, len(a))
	for i, v := range a {
		strList[i] = strconv.Itoa(v)
	}
	fmt.Printf("%s[%d] = np.array([%s])\n", name, idx, strings.Join(strList, ", "))
}
