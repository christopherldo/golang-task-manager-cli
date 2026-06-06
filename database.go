package main

import (
	"encoding/json"
	"log"
	"os"
)

func createDatabaseFile() *os.File {
	file, err := os.OpenFile(DB_URL, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Failed to open/create file: %s", err)
	}

	tasks := getAllTasksFromDatabase()

	if len(tasks) == 0 {
		writeDatabase([]Task{})
	}

	return file
}

func readDatabase() []byte {
	fileBytes, err := os.ReadFile(DB_URL)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	return fileBytes
}

func writeDatabase(tasks []Task) {
	fileData, err := json.MarshalIndent(tasks, "", "  ")

	if err != nil {
		log.Fatalf("Failed enconding JSON: %s", err)
	}

	err = os.WriteFile(DB_URL, fileData, 0644)

	if err != nil {
		log.Fatalf("Failed to write to database file")
	}
}
