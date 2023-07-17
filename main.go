package main

import (
	"fmt"
	mapReduce "here/localMapReduce"
	"os"
)

func main() {
	fmt.Println(os.Args[1])
	if os.Args[1] == "sequential" {
		// fmt.Println("here")
		SerialWordCounter()
	} else if os.Args[1] == "concurrent" {
		ConcCoordinator()
	} else {
		mapReduce.Coordinator()
	}

}
