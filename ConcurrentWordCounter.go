package main

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
	"time"
)

var workerNum int = 4
var wordSum int = 0
var wordMap = make(map[string]uint16)
var mutex = &sync.Mutex{} // Mutex to synchronize access to wordMap

func worker(id int, jobs chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := range jobs {
		fmt.Println("worker", id, "started new job")
		words := strings.Fields(j)
		countMap := make(map[string]uint16)

		for _, word := range words {
			word = strings.Trim(word, "!,.?!:;()'\"")
			// fmt.Println(word)
			countMap[word]++
		}

		mutex.Lock()
		for word, count := range countMap {
			// fmt.Println(word, count)
			wordMap[word] += count
		}
		mutex.Unlock()

		wordSum += len(words)
	}
}

func processDataConcurrently(buffer []byte, reader *bufio.Reader, jobs chan string) {
	var wg sync.WaitGroup

	for w := 1; w <= workerNum; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil {
			break
		}
		chunk := string(buffer[:bytesRead])

		jobs <- chunk
	}

	close(jobs)
	wg.Wait()
}

func ConcCoordinator() {
	fmt.Println("Inside ConcCoordinator")
	filePath := "./bible.txt"
	file := readFunc(filePath)
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 64*1024)

	jobs := make(chan string)

	processDataConcurrently(buffer, reader, jobs)

	fmt.Println(wordSum)
	time.Sleep(time.Second)
	fmt.Println(wordSum)
	fmt.Println(wordMap)
	defer fmt.Println("End of ConcCoordinator")
}
