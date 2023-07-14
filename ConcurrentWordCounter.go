package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

var workerNum int = 4
var wordSum int = 0

func worker(id int, jobs chan string) {
	for j := range jobs { //j is the 64kb chunk that is being read
		fmt.Println("worker", id, "started new job")
		// <-jobs
		words := strings.Fields(j)
		wordSum += len(words)
		fmt.Println("Num of chars in chunk", len(j))
		fmt.Println("Num of words in chunk", len(words))

		// fmt.Println("worker", id, "finished", j)

	}
}

func processDataConcurrently(buffer []byte, reader *bufio.Reader, jobs chan string) {
	// Create workers to receive chunks/jobs from job channel
	for w := 1; w <= workerNum; w++ {
		go worker(w, jobs)
	}

	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil {
			break
		}
		chunk := string(buffer[:bytesRead])

		//give Chunks to a job channel
		jobs <- chunk
		// break // <- delete me when done testing
	}
	close(jobs)
}

func Coordinator() {
	// Opening file
	fmt.Println("Inside Coordinator")
	filePath := "./bible.txt"
	file := readFunc(filePath)
	defer file.Close()

	// Create a buffered reader for efficient reading
	reader := bufio.NewReader(file)
	buffer := make([]byte, 64*1024)

	// Synchronization and Channels
	jobs := make(chan string)
	results := make(chan string)
	fmt.Println("placeholder: ", results)

	// Reading through file
	processDataConcurrently(buffer, reader, jobs)

	// fmt.Println(wordCount)

	fmt.Println(wordSum)
	time.Sleep(time.Second)
	fmt.Println(wordSum)
	defer fmt.Println("End of Coordinator")
}
