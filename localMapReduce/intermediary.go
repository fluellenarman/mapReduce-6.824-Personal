package mapReduce

import (
	"log"
	"sync"
)

var wgGlobalIntermediary sync.WaitGroup

// Each pair is assigned a reducer
var pairMap = make(map[string]int)

var curReducerProcess = 1
var ReducerJobs = make(chan keyValuePair)
var AvailableReducerProcesses = [2]bool{true, true}

var intermediaryLock = &sync.Mutex{}

func assignPairToReducer(pair keyValuePair) {
	_, ok := pairMap[pair.key]
	if !ok {
		// log.Println(pair)
		pairMap[pair.key] = curReducerProcess
		// log.Println(pairMap)

		if curReducerProcess < len(AvailableReducerProcesses) {
			curReducerProcess++
		} else {
			curReducerProcess = 1
		}
	}
}

func intermediaryWorker(id int, processId int, intermediaryJobs chan keyValuePair, wgIntermediary *sync.WaitGroup) {

	log.Println("Intermediary process", processId, ": worker", id, "started")

	for j := range intermediaryJobs {
		j = j
		intermediaryLock.Lock()
		assignPairToReducer(j)
		intermediaryLock.Unlock()
		// log.Println("Intermediary process", processId, ": worker", id, "processed", j.key)
	}

	log.Println("Intermediary process", processId, ": worker", id, "Finished")
	wgIntermediary.Done()
}

func IntermediaryProcess(id int, intermediaryJobs chan keyValuePair) {
	wgGlobalIntermediary.Add(1)
	var wgIntermediary sync.WaitGroup
	log.Println("Intermediary process", id, ": started")

	for w := 1; w <= 4; w++ {
		wgIntermediary.Add(1)
		go intermediaryWorker(w, id, intermediaryJobs, &wgIntermediary)
	}

	wgIntermediary.Wait()
	AvailableIntermediaryProcesses[id] = true
	// log.Println(pairMap)
	log.Println("Intermediary process", id, ": finished")
	wgGlobalIntermediary.Done()
}
