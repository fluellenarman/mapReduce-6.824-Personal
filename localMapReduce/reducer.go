package mapReduce

import (
	"log"
	"math/rand"
	"sync"
)

var wgGlobalReducer sync.WaitGroup
var finalOutputMap = make(map[string]uint16)
var finalWordCount = 0

var mutex = &sync.Mutex{} // Mutex to synchronize access to finalOutputMap

// merge finalOutputMap with workerOutputMap and update finalWordCount
func processFinalOutput(workerOutputMap map[string]uint16, wordCounter *int) {
	for key, value := range workerOutputMap {
		finalOutputMap[key] += value
	}

	finalWordCount += *wordCounter
	*wordCounter = 0
}

func reducerWorker(id int, processId int, reducerJobs chan keyValuePair, wgReducer *sync.WaitGroup) {
	log.Println("Reducer process", processId, ": worker", id, "started")
	workerOutputMap := make(map[string]uint16)
	counter := 0
	wordCounter := 0

	for j := range reducerJobs {
		workerOutputMap[j.key] += j.value
		wordCounter++

		// mutex.Lock()
		// processFinalOutput(workerOutputMap, &wordCounter)
		// mutex.Unlock()
		// workerOutputMap = make(map[string]uint16)

		if counter >= 100 {
			if mutex.TryLock() {
				processFinalOutput(workerOutputMap, &wordCounter)
				mutex.Unlock()
				workerOutputMap = make(map[string]uint16)
				counter = 0
			}
		} else {
			counter += rand.Intn(3)
		}
	}

	if len(workerOutputMap) != 0 {
		log.Println("Reducer process", processId, ": worker", id, "sequentializing")
		mutex.Lock()
		processFinalOutput(workerOutputMap, &wordCounter)
		mutex.Unlock()
		workerOutputMap = make(map[string]uint16)
		counter = 0
	}

	// log.Println("Reducer process", processId, ": worker", id, "processed", len(workerOutputMap))
	log.Println("Reducer process", processId, ": worker", id, "finished")
	wgReducer.Done()
}

func ReducerProcess(id int, reducerJobs chan keyValuePair) {
	log.Println("Reducer process", id, "started")

	wgGlobalReducer.Add(1)
	var wgReducer sync.WaitGroup

	for w := 1; w <= 4; w++ {
		wgReducer.Add(1)
		go reducerWorker(w, id, reducerJobs, &wgReducer)
	}
	wgReducer.Wait()

	log.Println("Reducer process", id, "finished")
	wgGlobalReducer.Done()
}
