package mapReduce

import (
	"log"
	"sync"
)

var wgGlobalIntermediary sync.WaitGroup

var AvailableReducerProcesses = [2]bool{true, true}

var reducerChannels [2]chan keyValuePair

func activateReducerProcess(reducerChannels *[2]chan keyValuePair) {
	for i, available := range AvailableReducerProcesses {
		if available {
			AvailableReducerProcesses[i] = false
			go ReducerProcess(i, reducerChannels[i])
			break
		} else {
			// log.Println("Coordinator: Reducer process", i, "is not available")
		}
	}
}

func assignPairToReducer(pair keyValuePair, pairMap map[string]int, currentReducer *int) {
	_, ok := pairMap[pair.key]
	if !ok {
		pairMap[pair.key] = *currentReducer
	}

	if *currentReducer >= len(AvailableReducerProcesses)-1 {
		*currentReducer = 0
	} else {
		*currentReducer++
	}
}

func intermediaryWorker(id int, processId int, intermediaryJobs chan keyValuePair, wgIntermediary *sync.WaitGroup) {

	log.Println("Intermediary process", processId, ": worker", id, "started")

	// Each pair is assigned a reducer
	// Each worker has their own pairmap
	var pairMap = make(map[string]int)
	currentReducer := 0

	for j := range intermediaryJobs {
		activateReducerProcess(&reducerChannels)
		assignPairToReducer(j, pairMap, &currentReducer)
		reducerChannels[pairMap[j.key]] <- j
		// log.Println(pairMap[j.key])
	}

	log.Println("Intermediary process", processId, ": worker", id, "Finished")
	wgIntermediary.Done()
}

func IntermediaryProcess(id int, intermediaryJobs chan keyValuePair) {
	log.Println("Intermediary process", id, ": started")

	wgGlobalIntermediary.Add(1)
	var wgIntermediary sync.WaitGroup

	// Create workers
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
