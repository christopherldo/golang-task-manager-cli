package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"chrisldo.com/todo-cli/internal/todo/models"
)

var DbUrl = "db.json"
var dbMutex sync.RWMutex

func EnsureDataBaseExists() error {
	_, err := os.Stat(DbUrl)

	if os.IsNotExist(err) {
		err = WriteDatabase([]models.Task{})

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

func ReadDatabase() ([]byte, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	fileBytes, err := os.ReadFile(DbUrl)

	if err != nil {
		return nil, fmt.Errorf("Error reading file: %w", err)
	}

	return fileBytes, nil
}

func WriteDatabase(tasks []models.Task) error {
	fileData, err := json.MarshalIndent(tasks, "", "  ")

	if err != nil {
		return fmt.Errorf("Failed enconding JSON: %w", err)
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	err = os.WriteFile(DbUrl, fileData, 0644)

	if err != nil {
		return fmt.Errorf("Failed to write to database file: %w", err)
	}

	return nil
}
