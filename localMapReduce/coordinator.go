package mapReduce

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

var availableMapperProcesses = [2]bool{true, true}

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
		activateMapperProcess(mapperJobs)

		// channels represent sending and receiving data
		mapperJobs <- chunk
		log.Println("Coordinator: Chunk sent to Mapper process")
		// break //<- For testing
	}
	wg.Done()
	close(mapperJobs)
}

func Coordinator() {
	// log.SetFlags(log.Lshortfile)
	log.SetFlags(log.Lmicroseconds)
	log.Println("In coordinator")

	// Open file
	file, err := os.Open("./bible.txt")
	checkNilErr(err)
	defer file.Close()

	// Initialize reader and buffer
	reader := bufio.NewReader(file)
	buffer := make([]byte, 64*1024) // 64kb chunks
	// initialize channel
	mapperJobs := make(chan string, 100)

	var wg sync.WaitGroup
	wg.Add(1)

	readData(reader, buffer, mapperJobs, &wg)

	wg.Wait()
	log.Println("Coordinator: Finished reading file")

	log.Println("Coordinator: Waiting for mapper processes to finish")
	wgGlobalMapper.Wait()
	log.Println("Coordinator: mapper proccesses finished")

	close(IntermediaryJobs)
	log.Println("Coordinator: Closed intermediary channel")

	log.Println("Coordinator: Waiting for intermediary processes to finish")
	wgGlobalIntermediary.Wait()

	log.Println("Coordinator: Program finished")
}
