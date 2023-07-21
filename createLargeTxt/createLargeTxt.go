package createlargetxt

import (
	"io"
	"log"
	"os"
)

func CreateLargeTxt() {
	filePath := "./bible.txt"
	outputFilePath := "./bible100.txt"

	file, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Panic(err)
		return
	}

	// Open a new file in write mode
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Panic(err)
		return
	}
	defer outputFile.Close()

	// Write the content of "bible.txt" to the new file 100 times
	for i := 0; i < 100; i++ {
		if _, err := outputFile.Write(content); err != nil {
			log.Panic(err)
			return
		}
	}

	log.Println("New file 'bible100.txt' has been created")
}
