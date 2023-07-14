package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args[1])
	if os.Args[1] == "sequential" {
		// fmt.Println("here")
		SerialWordCounter()
	} else {
		Coordinator()
	}

}
