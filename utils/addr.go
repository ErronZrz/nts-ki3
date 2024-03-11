package utils

import (
	"errors"
)

var (
	toAdd = []int{
		1, 1, 1, 3, 1, 5, 3, 3, 1, 9, 7,
		5, 3, 17, 27, 3, 1, 29, 3, 21, 7, 17,
		15, 9, 43, 35, 15, 29, 3, 11, 3, 11, 15,
	}
	seeds = []int{
		-1, 2, 2, 2, 3, 2, 2, 2, 3, 3, 14,
		2, 2, 7, 3, 2, 3, 17, 2, 2, 5, 47,
		3, 3, 2, 2, 3, 5, 2, 3, 2, 2, 3,
	}
)

func CidrPow(cidr string) (int, error) {
	n := len(cidr)
	pow := 32
	val := cidr[n-1]
	if val < 0x30 || val > 0x39 {
		return -1, errors.New("invalid CIDR address")
	}
	pow -= int(val - 0x30)
	val = cidr[n-2]
	if val == 0x2F {
		return pow, nil
	}
	if cidr[n-3] != 0x2F || val < 0x31 || val > 0x33 {
		return -1, errors.New("invalid CIDR address")
	}
	pow -= 10 * int(val-0x30)
	if pow < 0 {
		return -1, errors.New("invalid CIDR address")
	}
	return pow, nil
}

func SmallestPrime(pow int) (int, error) {
	if pow < 0 || pow > 32 {
		return -1, errors.New("unsupported number")
	}
	return (1 << pow) + toAdd[pow], nil
}

func GetSeed(pow int) (int, error) {
	if pow < 0 || pow > 32 {
		return -1, errors.New("unsupported number")
	}
	return seeds[pow], nil
}

func FindSeed(root int) int {
	if root <= 3 {
		return root - 1
	}
	seed := 2
	for {
		val := (seed * seed) % root
		times := 1
		for val != seed {
			times++
			val = (val * seed) % root
		}
		if times+1 == root {
			break
		}
		seed++
	}
	return seed
}
