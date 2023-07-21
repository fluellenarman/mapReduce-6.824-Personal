package mapReduce

import (
	"log"
	"sync"
)

var wgGlobalReducer sync.WaitGroup

func reducerWorker(id int, processId int, reducerJobs chan keyValuePair, wgReducer *sync.WaitGroup) {
	log.Println("Reducer process", processId, ": worker", id, "started")

	for j := range reducerJobs {
		// log.Println(j)
		j = j
	}

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
