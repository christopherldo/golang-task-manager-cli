package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"chrisldo.com/todo-cli/internal/todo/models"
)

type FileStore struct {
	mu       sync.RWMutex
	filePath string
}

func NewFileStore(path string) (*FileStore, error) {
	store := &FileStore{
		filePath: path,
	}

	err := store.ensureDataBaseExists()

	if err != nil {
		return nil, fmt.Errorf("Failed to initialize database: %w", err)
	}

	return store, nil
}

func (fs *FileStore) ensureDataBaseExists() error {
	_, err := os.Stat(fs.filePath)

	if os.IsNotExist(err) {
		err = fs.WriteDatabase([]models.Task{})

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

func (fs *FileStore) ReadDatabase() ([]byte, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	fileBytes, err := os.ReadFile(fs.filePath)

	if err != nil {
		return nil, fmt.Errorf("Error reading file: %w", err)
	}

	return fileBytes, nil
}

func (fs *FileStore) WriteDatabase(tasks []models.Task) error {
	fileData, err := json.MarshalIndent(tasks, "", "  ")

	if err != nil {
		return fmt.Errorf("Failed enconding JSON: %w", err)
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	err = os.WriteFile(fs.filePath, fileData, 0644)

	if err != nil {
		return fmt.Errorf("Failed to write to database file: %w", err)
	}

	return nil
}
