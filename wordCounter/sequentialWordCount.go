package wordCounter

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

func avgArray(input []float32) float32 {
	i := float32(0)
	total := float32(0)
	for _, val := range input {
		total += val
		i += 1
	}
	return total / i
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
			// log.Println(words[i])
		}
	}
	return wordCount, wordMap
}

func SequentialWordCounter() {
	// Declaring Variables
	wordCount := 0
	wordMap := map[string]uint16{}

	// Open the file
	filePath := "./bible.txt"
	file := readFunc(filePath)
	defer file.Close()

	// Create a buffered reader for efficient reading
	reader := bufio.NewReader(file)
	buffer := make([]byte, 64*1024) // 64KB buffer size

	// Initializing timers
	localTimeArray := []float32{}
	startTotal := time.Now()
	elapsedTotal := time.Since(startTotal)

	for i := 0; i < 1; i++ {
		startLocal := time.Now()
		wordCount, wordMap = processData(buffer, reader)
		elapsedLocal := time.Since(startLocal)
		elapsedTotal = time.Since(startTotal)
		localTimeArray = append(localTimeArray, float32(elapsedLocal.Milliseconds()))

		fmt.Println(elapsedTotal, elapsedLocal)

		file = readFunc(filePath)
		reader = bufio.NewReader(file)
	}
	fmt.Println(wordCount)

	elapsedTotal = time.Since(startTotal)

	fmt.Println(wordMap)

	fmt.Println("Total word count:", wordCount)
	fmt.Printf("Total Execution time: %s\n", elapsedTotal)
	fmt.Printf("Avg Execution time: %vms\n", avgArray(localTimeArray))

}
