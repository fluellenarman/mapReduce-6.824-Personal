package mapReduce

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var availableMapperProcesses = [2]bool{true, true}

func writeOutputToFile(finalOutputMap map[string]uint16) {
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	writer := bufio.NewWriter(outputFile)
	for key, value := range finalOutputMap {
		_, err := fmt.Fprintf(writer, "%s: %v\n", key, value)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("Error flushing data:", err)
			return
		}
	}

	outputFile.Close()
}

func initializeReducerChannels(channels *[2]chan keyValuePair) {
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan keyValuePair, 100)
	}
}

func checkNilErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func activateMapperProcess(mapperJobs chan string) {
	for i, available := range availableMapperProcesses {
		if available {
			availableMapperProcesses[i] = false
			go MapperProcess(i, mapperJobs)
			break
		} else {
			// log.Println("Coordinator: Mapper process", i, "is not available")
		}
	}
}

func readData(reader *bufio.Reader, buffer []byte, mapperJobs chan string, wg *sync.WaitGroup) {
	log.Println("Coordinator: Started reading data")
	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil {
			break
		}
		chunk := string(buffer[:bytesRead])
		// activate mapperProcess before sending data
		for i, available := range availableMapperProcesses {
			if available {
				go MapperProcess(i, mapperJobs)
				availableMapperProcesses[i] = false
			}
		}

		// channels represent sending and receiving data
		mapperJobs <- chunk
		log.Println("Coordinator: Chunk sent to Mapper process")
		// break //<- For testing
	}
	wg.Done()
}

func Coordinator() {
	var wg sync.WaitGroup
	wg.Add(1)

	log.SetFlags(log.Lshortfile)
	// log.SetFlags(log.Lmicroseconds)

	startTime := time.Now()

	// Initialize channels
	initializeReducerChannels(&reducerChannels)

	// Open file
	file, err := os.Open("./bible.txt")
	checkNilErr(err)
	defer file.Close()

	// Initialize reader and buffer
	reader := bufio.NewReader(file)
	buffer := make([]byte, 64*1024) // 64kb chunks

	// initialize channel
	mapperJobs := make(chan string, 100)

	activateMapperProcess(mapperJobs)

	readData(reader, buffer, mapperJobs, &wg)

	// Waiting for proccesses to finish and closing channels
	wg.Wait()
	close(mapperJobs)
	log.Println("Coordinator: Closed mapper channel")
	log.Println("Coordinator: Finished reading file")

	log.Println("Coordinator: Waiting for mapper processes to finish")
	wgGlobalMapper.Wait()
	log.Println("Coordinator: mapper proccesses finished")
	close(IntermediaryJobs)
	log.Println("Coordinator: Closed intermediary channel")

	log.Println("Coordinator: Waiting for intermediary processes to finish")
	wgGlobalIntermediary.Wait()
	log.Println("Coordinator: intermediary processes finished")
	for i := 0; i < len(reducerChannels); i++ {
		close(reducerChannels[i])
	}
	log.Println("Coordinator: Closed reducer channels")

	log.Println("Coordinator: Waiting for reducer processes to finish")
	wgGlobalReducer.Wait()
	log.Println("Coordinator: reducer processes finished")

	writeOutputToFile(finalOutputMap)

	log.Println("Coordinator: Program finished")

	elapsedTime := time.Since(startTime)
	fmt.Println("Time taken:", elapsedTime)

	// log.Println(finalOutputMap)
	log.Println(len(finalOutputMap))
	log.Println(finalWordCount)
}
