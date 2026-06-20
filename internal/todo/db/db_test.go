package db

import (
	"os"
	"testing"

	"chrisldo.com/todo-cli/internal/todo/models"
)

func TestWriteAndReadDatabase(t *testing.T) {
	DbUrl := "db_test.json"

	defer func() {
		os.Remove(DbUrl)
	}()

	dataBase, err := NewFileStore(DbUrl)

	err = dataBase.ensureDataBaseExists()

	if err != nil {
		t.Fatalf("Expected database file to exist")
	}

	mockTests := []models.Task{
		{ID: 1, Description: "Task de Teste", IsDone: false},
	}

	err = dataBase.WriteDatabase(mockTests)

	if err != nil {
		t.Fatalf("Expected to write without an error, but failed: %v", err)
	}

	bytes, err := dataBase.ReadDatabase()

	if err != nil {
		t.Fatalf("Expected to read the database file, but failed: %v", err)
	}

	if len(bytes) == 0 {
		t.Errorf("The database should not be empty")
	}
}
