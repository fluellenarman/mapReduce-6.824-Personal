package mapReduce

import (
	"log"
)

func intermediary(id int, intermediaryJobs chan map[string]uint16) {
	log.Println("Intermediary process", id, ": started")
	for j := range intermediaryJobs {
		log.Println("Intermediary process", id, ": received job")
		log.Println(len(j))
		log.Println("Intermediary process", id, ": finished processing job")
	}
	log.Println("Intermediary process", id, ": finished")
}
