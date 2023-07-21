package main

import (
	"fmt"
	mapReduce "here/localMapReduce"
	"os"
	"time"
)

func main() {
	fmt.Println(os.Args[1])
	startTime := time.Now()
	if os.Args[1] == "sequential" {
		// fmt.Println("here")
		SerialWordCounter()
	} else if os.Args[1] == "concurrent" {
		ConcCoordinator()
	} else if os.Args[1] == "mapreduce" {
		mapReduce.Coordinator()
	}
	elapsedTime := time.Since(startTime)
	fmt.Println("Time taken:", elapsedTime)
}
