package main

import (
	"encoding/json"
	"fmt"
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
		return nil, fmt.Errorf("Error reading file: %w", err)
	}

	return fileBytes, nil
}

func writeDatabase(tasks []Task) error {
	fileData, err := json.MarshalIndent(tasks, "", "  ")

	if err != nil {
		return fmt.Errorf("Failed enconding JSON: %w", err)
	}

	err = os.WriteFile(DB_URL, fileData, 0644)

	if err != nil {
		return fmt.Errorf("Failed to write to database file: %w", err)
	}

	return nil
}
