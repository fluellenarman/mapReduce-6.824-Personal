package mapReduce

import (
	"log"
	"strings"
	"sync"
)

func activateIntermediaryProcess(intermediaryJobs chan map[string]uint16) {
	for i, available := range AvailableIntermediaryProcesses {
		if available {
			AvailableIntermediaryProcesses[i] = false
			go intermediary(i, intermediaryJobs)
			break
		} else {
			log.Println("Coordinator: Intermediary process", i, "is not available")
		}
	}
}

func mapperWorker(id int, processId int, mapperJobs chan string, wgMapper *sync.WaitGroup) {
	log.Println("Mapper process", processId, ": worker", id, "started")
	// var wgLocal sync.WaitGroup

	for j := range mapperJobs {
		words := strings.Fields(j)
		wordMap := make(map[string]uint16)
		countLocal := 0

		for _, word := range words {
			countLocal++
			word = strings.Trim(word, "!,.?!:;()'\"")
			wordMap[word]++
			// log.Println("Mapper process", processId, ": worker", id, "processed", word)
		}
		log.Println("Mapper process", processId, ": worker", id, "processed", countLocal, "words")
		activateIntermediaryProcess(IntermediaryJobs)

		// Send wordMap and wordCount to intermediary/shuffle and sort
		IntermediaryJobs <- wordMap
		// log.Println(wordMap)
	}

	log.Println("Mapper process", processId, ": worker", id, "finished")
	wgMapper.Done()
}

// process chunks into key value pairs
// mapperJobs is a channel of chunks
func MapperProcess(id int, mapperJobs chan string) {
	var wgMapper sync.WaitGroup
	log.Println("Mapper process", id, ": started")
	for w := 1; w <= 4; w++ {
		wgMapper.Add(1)
		go mapperWorker(w, id, mapperJobs, &wgMapper)
	}

	// send to coordinator that job is done and process is available
	wgMapper.Wait()
	availableMapperProcesses[id] = true
	log.Println("Mapper process", id, ": finished")
	log.Println(availableMapperProcesses)
}
