package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"chrisldo.com/todo-cli/internal/todo/models"
)

type DatabaseSchema struct {
	LastID int           `json:"last_task_id"`
	Tasks  []models.Task `json:"tasks"`
}

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
		return nil, fmt.Errorf("initializing database: %w", err)
	}

	return store, nil
}

func (fs *FileStore) ensureDataBaseExists() error {
	_, err := os.Stat(fs.filePath)

	if os.IsNotExist(err) {
		err = fs.WriteDatabase(DatabaseSchema{
			LastID: 0,
			Tasks:  []models.Task{},
		})

		if err != nil {
			return fmt.Errorf("initializing database: %w", err)
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("ensuring database exists: %w", err)
	}

	return nil
}

func (fs *FileStore) ReadDatabase() (DatabaseSchema, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var schema DatabaseSchema

	fileBytes, err := os.ReadFile(fs.filePath)

	if err != nil {
		return schema, fmt.Errorf("reading file: %w", err)
	}

	if len(fileBytes) == 0 {
		return schema, nil
	}

	if err := json.Unmarshal(fileBytes, &schema); err != nil {
		return schema, fmt.Errorf("parsing database JSON: %w", err)
	}

	return schema, nil
}

func (fs *FileStore) WriteDatabase(schema DatabaseSchema) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fileBytes, err := json.Marshal(schema)

	if err != nil {
		return fmt.Errorf("parsing JSON to bytes: %w", err)
	}

	err = os.WriteFile(fs.filePath, fileBytes, 0644)

	if err != nil {
		return fmt.Errorf("writing to database file: %w", err)
	}

	return nil
}
