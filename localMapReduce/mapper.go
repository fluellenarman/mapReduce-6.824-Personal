package mapReduce

import (
	"log"
	"strings"
	"sync"
)

type keyValuePair struct {
	key   string
	value uint16
}

var wgGlobalMapper sync.WaitGroup

var AvailableIntermediaryProcesses = [2]bool{true, true}
var IntermediaryJobs = make(chan keyValuePair, 100)

func activateIntermediaryProcess(intermediaryJobs chan keyValuePair) {
	for i, available := range AvailableIntermediaryProcesses {
		if available {
			AvailableIntermediaryProcesses[i] = false
			go IntermediaryProcess(i, intermediaryJobs)
			break
		} else {
			// log.Println("Mapper: Intermediary process", i, "is not available")
		}
	}
}

func mapperWorker(id int, processId int, mapperJobs chan string, wgMapper *sync.WaitGroup) {
	log.Println("Mapper process", processId, ": worker", id, "started")

	for j := range mapperJobs {
		words := strings.Fields(j)
		countLocal := 0

		activateIntermediaryProcess(IntermediaryJobs)

		for _, word := range words {
			countLocal++
			word = strings.Trim(word, "!,.?!:;()'\"")
			pair := keyValuePair{word, 1}
			// pair = pair // for testing

			// Send key value pair to intermediary process
			IntermediaryJobs <- pair
		}
		log.Println("Mapper process", processId, ": worker", id, "processed", countLocal, "words")
	}

	log.Println("Mapper process", processId, ": worker", id, "finished")
	wgMapper.Done()
}

// process chunks into key value pairs
// mapperJobs is a channel of chunks
func MapperProcess(id int, mapperJobs chan string) {
	var wgMapper sync.WaitGroup
	wgGlobalMapper.Add(1)
	log.Println("Mapper process", id, ": started")
	for w := 1; w <= 4; w++ {
		wgMapper.Add(1)
		go mapperWorker(w, id, mapperJobs, &wgMapper)
	}

	// send to coordinator that job is done and process is available
	wgMapper.Wait()
	availableMapperProcesses[id] = true
	log.Println("Mapper process", id, ": finished")

	wgGlobalMapper.Done()
}
