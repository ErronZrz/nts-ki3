package main

import (
	"active/offset"
	"flag"
	"fmt"
	"time"
)

func main() {
	var round, interval, roundInterval int
	var inputPath, outputPath string
	flag.IntVar(&round, "round", 1, "Number of rounds to run")
	flag.StringVar(&inputPath, "input", "", "Input text file path")
	flag.StringVar(&outputPath, "output", "", "Output text file path")
	flag.IntVar(&interval, "interval", 1000, "Interval between tasks in milliseconds")
	flag.IntVar(&roundInterval, "roundInterval", 60, "Interval between rounds in seconds")
	flag.Parse()

	for i := 0; i < round; i++ {
		fmt.Printf("Round %d started.\n", i)
		err := offset.CalculateOffsetsAsync(inputPath, outputPath, interval)
		if err != nil {
			fmt.Println(err)
		}
		if i < round-1 {
			fmt.Printf("Round %d finished. Now wait for %d seconds.\n", i, roundInterval)
			time.Sleep(time.Duration(roundInterval) * time.Second)
		}
	}
}
