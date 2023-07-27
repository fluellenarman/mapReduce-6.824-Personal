package main

import (
	"fmt"
	helperTxt "here/createLargeTxt"
	mapReduce "here/localMapReduce"
	wordCounter "here/wordCounter"
	"os"
)

func main() {
	fmt.Println(os.Args[1])
	// startTime := time.Now()
	if os.Args[1] == "sequential" {
		wordCounter.SequentialWordCounter()
	} else if os.Args[1] == "concurrent" {
		wordCounter.ConcCoordinator()
	} else if os.Args[1] == "mapreduce" {
		mapReduce.Coordinator()
	} else if os.Args[1] == "createLargeTxt" {
		helperTxt.CreateLargeTxt()
	}
	// elapsedTime := time.Since(startTime)
	// fmt.Println("Time taken:", elapsedTime)
}
