package mapReduce

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

var availableMapperProcesses = [2]bool{true, true}
var AvailableIntermediaryProcesses = [2]bool{true, true}

var IntermediaryJobs = make(chan map[string]uint16)

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
			log.Println("Coordinator: Mapper process", i, "is not available")
		}
	}
}

func readData(reader *bufio.Reader, buffer []byte, mapperJobs chan string, wg *sync.WaitGroup) {
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
	log.Println("Coordinator: Finished reading data")
	close(mapperJobs)
}

func Coordinator() {
	log.SetFlags(log.Lshortfile)
	log.Println("In coordinator")

	// Open file
	file, err := os.Open("./bible.txt")
	checkNilErr(err)
	defer file.Close()

	// Initialize reader and buffer
	reader := bufio.NewReader(file)
	buffer := make([]byte, 64*1024) // 64kb chunks
	// initialize channel
	mapperJobs := make(chan string)

	var wg sync.WaitGroup
	wg.Add(1)

	readData(reader, buffer, mapperJobs, &wg)

	wg.Wait()
	// time.Sleep(time.Second * 2) //<- For testing , DELETE LATER
}
