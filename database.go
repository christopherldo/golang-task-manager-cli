package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	DB_URL = "db.json"
)

func ensureDataBaseExists() {
	_, err := os.Stat(DB_URL)

	if os.IsNotExist(err) {
		writeDatabase([]Task{})
	}
}

func readDatabase() ([]byte, error) {
	fileBytes, err := os.ReadFile(DB_URL)

	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return fileBytes, err
	}

	return fileBytes, err
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
