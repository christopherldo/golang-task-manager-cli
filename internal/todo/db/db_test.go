package db

import (
	"os"
	"testing"

	"chrisldo.com/todo-cli/internal/todo/models"
)

func TestWriteAndReadDatabase(t *testing.T) {
	originalDb := DbUrl
	DbUrl = "test_db.json"

	defer func() {
		os.Remove(DbUrl)
		DbUrl = originalDb
	}()

	err := EnsureDataBaseExists()

	if err != nil {
		t.Fatalf("Expected database file to exist")
	}

	mockTests := []models.Task{
		{ID: 1, Description: "Task de Teste", IsDone: false},
	}

	err = WriteDatabase(mockTests)

	if err != nil {
		t.Fatalf("Expected to write without an error, but failed: %v", err)
	}

	bytes, err := ReadDatabase()

	if err != nil {
		t.Fatalf("Expected to read the database file, but failed: %v", err)
	}

	if len(bytes) == 0 {
		t.Errorf("The database should not be empty")
	}
}
