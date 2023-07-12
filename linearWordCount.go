package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func checkNilErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func readFunc(filePath string) *os.File {
	file, err := os.Open(filePath)
	checkNilErr(err)
	return file
	// defer file.Close()
}

func processData(buffer []byte, reader *bufio.Reader) (int, map[string]uint16) {
	wordCount := 0
	wordMap := map[string]uint16{}

	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil {
			break
		}

		// Count the words in the chunk
		chunk := string(buffer[:bytesRead])
		words := strings.Fields(chunk)
		wordCount += len(words)
		for i := 0; i < len(words); i++ {
			// fmt.Println(words[i])
			words[i] = strings.Trim(words[i], "!,.?!:;()'\"")
			wordMap[words[i]] += 1
		}
	}
	return wordCount, wordMap
}

func main() {
	// Open the file
	filePath := "./bible.txt"
	file := readFunc(filePath)
	defer file.Close()

	// Create a buffered reader for efficient reading
	reader := bufio.NewReader(file)
	buffer := make([]byte, 64*1024) // 64KB buffer size

	startTotal := time.Now()
	elapsedTotal := time.Since(startTotal)
	wordCount := 0
	wordMap := map[string]uint16{}
	for i := 0; i < 100; i++ {
		startLocal := time.Now()
		wordCount, wordMap = processData(buffer, reader)
		elapsedLocal := time.Since(startLocal)
		elapsedTotal = time.Since(startTotal)

		fmt.Println(elapsedTotal, elapsedLocal)

		file = readFunc(filePath)
		reader = bufio.NewReader(file)
	}
	fmt.Println(wordCount)

	elapsedTotal = time.Since(startTotal)

	fmt.Println(wordMap)

	fmt.Println("Total word count:", wordCount)
	fmt.Printf("Total Execution time: %s\n", elapsedTotal)
	fmt.Printf("Avg Execution time: %s\n", elapsedTotal/100)

}
