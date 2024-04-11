package utils

import (
	"fmt"
	"testing"
)

func TestValueCount(t *testing.T) {
	path := "D:\\Desktop\\TMP\\ntpdata\\2024-04-11_ntske_1.txt"
	columns := []int{6, 7, 8, 9, 10, 11}
	countMaps, err := ValueCount(path, columns)
	if err != nil {
		t.Error(err)
	}
	for _, countMap := range countMaps {
		fmt.Println(countMap)
	}
}
