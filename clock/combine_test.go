package clock

import (
	"fmt"
	"testing"
)

func TestCombineAlgorithm(t *testing.T) {
	maxSurvivors := 40
	result, sj := ClusterAlgorithm(data, maxSurvivors)
	for _, p := range result {
		fmt.Printf("%+v\n", p)
	}

	values := CombineAlgorithm(result, sj)
	fmt.Printf("System Peer:\n%+v\n", result[0])
	fmt.Printf("%+v\n", values)
}
