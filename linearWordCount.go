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

func main() {
	// Open the file
	file, err := os.Open("./bible.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	start := time.Now()

	// Create a buffered reader for efficient reading
	reader := bufio.NewReader(file)

	wordCount := 0
	var words []string

	// charsToCut := "!,.?!:;'\""

	wordMap := map[string]uint16{}

	// Read the file chunk by chunk until EOF
	buffer := make([]byte, 64*1024) // 64KB buffer size
	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil {
			break
		}

		// Count the words in the chunk
		chunk := string(buffer[:bytesRead])
		words = strings.Fields(chunk)
		wordCount += len(words)
		for i := 0; i < len(words); i++ {
			// fmt.Println(words[i])
			words[i] = strings.Trim(words[i], "!,.?!:;()'\"")
			wordMap[words[i]] += 1
		}
	}
	// fmt.Println(words[0])
	elapsed := time.Since(start)

	fmt.Println(len(words))
	fmt.Println(wordMap)

	fmt.Println("Total word count:", wordCount)
	fmt.Printf("Execution time: %s\n", elapsed)

}
