package repository

import (
	"os"
	"testing"

	"chrisldo.com/todo-cli/internal/todo/db"
	"chrisldo.com/todo-cli/internal/todo/models"
)

func resetState() {
	cacheMutex.Lock()
	cachedTasks = []models.Task{}
	cachedLastTaskId = nil
	cacheMutex.Unlock()
}

func setupTestEnvironment() func() {
	originalDB := db.DbUrl
	db.DbUrl = "test_repository_db.json"
	resetState()

	return func() {
		db.DbUrl = originalDB
		os.Remove("test_repository_db.json")
	}
}

func TestAppendAndGetAllTasks(t *testing.T) {
	cleanup := setupTestEnvironment()
	defer cleanup()

	newTask := models.Task{ID: 1, Description: "Aprender Go Testing", IsDone: false}

	err := AppendTaskToDatabase(newTask)

	if err != nil {
		t.Fatalf("Failed trying to add task: %v", err)
	}

	tasks := GetAllTasksFromDatabase()

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task in memory, found %d", len(tasks))
	}

	if tasks[0].Description != "Aprender Go Testing" {
		t.Fatalf("Task description is not the same as expected")
	}
}

func TestUpdateTask(t *testing.T) {
	cleanup := setupTestEnvironment()
	defer cleanup()

	AppendTaskToDatabase(models.Task{ID: 1, Description: "Task Antiga", IsDone: false})

	updatedTask := models.Task{ID: 1, Description: "Task Nova", IsDone: true}

	err := UpdateTaskOnDatabase(updatedTask)

	if err != nil {
		t.Fatalf("Unexpected error while updating: %v", err)
	}

	taskInDb, _ := GetOneTaskFromDatabase(1)

	if taskInDb.Description != "Task Nova" || taskInDb.IsDone != true {
		t.Errorf("Task was not updated correctly")
	}
}

func TestDeleteTask(t *testing.T) {
	cleanup := setupTestEnvironment()
	defer cleanup()

	AppendTaskToDatabase(models.Task{ID: 10, Description: "Task para deletar", IsDone: false})

	err := DeleteTaskFromDatabase(10)

	if err != nil {
		t.Fatalf("Unexpected error while deleting task: %v", err)
	}

	_, err = GetOneTaskFromDatabase(10)

	if err == nil {
		t.Errorf("Expected an error while searching already deleted task, but task it still in memory")
	}
}

func TestMarkTaskAsDone(t *testing.T) {
	cleanup := setupTestEnvironment()
	defer cleanup()

	AppendTaskToDatabase(models.Task{ID: 1, Description: "Task para marcar como completada", IsDone: false})

	err := MarkTaskAsDone(1)

	if err != nil {
		t.Fatalf("Unexpected error while marking task as done: %v", err)
	}

	taskInDb, _ := GetOneTaskFromDatabase(1)

	if taskInDb.IsDone != true {
		t.Errorf("Task was not set as done correctly")
	}
}

func TestGetLastTaskId(t *testing.T) {
	cleanup := setupTestEnvironment()
	defer cleanup()

	AppendTaskToDatabase(models.Task{})

	lastTaskId := GetLastTaskId()

	if lastTaskId != 0 {
		t.Errorf("Expected the last task id to be exactly 0 when there is no tasks added yet.")
	}

	AppendTaskToDatabase(models.Task{ID: 1, Description: "Task 1", IsDone: false})
	AppendTaskToDatabase(models.Task{ID: 1, Description: "Task 2", IsDone: false})
	// ...
	AppendTaskToDatabase(models.Task{ID: 100, Description: "Task 100", IsDone: false})

	lastTaskId = GetLastTaskId()

	if lastTaskId != 100 {
		t.Errorf("Expected the last task id to be exactly 100")
	}
}

func TestLoadDatabaseToMemory(t *testing.T) {
	cleanup := setupTestEnvironment()
	defer cleanup()

	db.WriteDatabase([]models.Task{
		{ID: 1, Description: "Task 1", IsDone: false},
	})

	err := LoadDatabaseToMemory()

	if err != nil {
		t.Fatalf("Unexpected error while loading the database file to the memory: %v", err)
	}

	if len(cachedTasks) != 1 {
		t.Errorf("Expected exactly 1 task loaded in memory, found %d", len(cachedTasks))
	}
}
