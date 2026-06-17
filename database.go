package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	DB_URL = "db.json"
)

var dbMutex sync.RWMutex

func ensureDataBaseExists() error {
	_, err := os.Stat(DB_URL)

	if os.IsNotExist(err) {
		err = writeDatabase([]Task{})

		if err != nil {
			return fmt.Errorf("Error initializing database: %w", err)
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("Error ensuring database exists: %w", err)
	}

	return nil
}

func readDatabase() ([]byte, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

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

	dbMutex.Lock()
	defer dbMutex.Unlock()

	err = os.WriteFile(DB_URL, fileData, 0644)

	if err != nil {
		return fmt.Errorf("Failed to write to database file: %w", err)
	}

	return nil
}
