package repository

import (
	"os"
	"testing"

	"chrisldo.com/todo-cli/internal/todo/db"
	"chrisldo.com/todo-cli/internal/todo/models"
)

func setupTestEnvironment(dbUrl string) (func(), Store) {
	store, _ := db.NewFileStore(dbUrl)

	return func() {
		os.Remove(dbUrl)
	}, store
}

func TestAppendAndGetAllTasks(t *testing.T) {
	cleanup, store := setupTestEnvironment("db_TestAppendAndGetAllTasks_test.json")
	defer cleanup()

	repo := NewTaskRepository(store)

	newTask := models.Task{ID: 1, Description: "Aprender Go Testing", IsDone: false}

	err := repo.AppendTaskToDatabase(newTask)

	if err != nil {
		t.Fatalf("Failed trying to add task: %v", err)
	}

	tasks := repo.GetAllTasksFromDatabase()

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task in memory, found %d", len(tasks))
	}

	if tasks[0].Description != "Aprender Go Testing" {
		t.Fatalf("Task description is not the same as expected")
	}
}

func TestUpdateTask(t *testing.T) {
	cleanup, store := setupTestEnvironment("db_TestUpdateTask_test.json")
	defer cleanup()

	repo := NewTaskRepository(store)

	repo.AppendTaskToDatabase(models.Task{ID: 1, Description: "Task Antiga", IsDone: false})

	updatedTask := models.Task{ID: 1, Description: "Task Nova", IsDone: true}

	err := repo.UpdateTaskOnDatabase(updatedTask)

	if err != nil {
		t.Fatalf("Unexpected error while updating: %v", err)
	}

	taskInDb, _ := repo.GetOneTaskFromDatabase(1)

	if taskInDb.Description != "Task Nova" || taskInDb.IsDone != true {
		t.Errorf("Task was not updated correctly")
	}
}

func TestDeleteTask(t *testing.T) {
	cleanup, store := setupTestEnvironment("db_TestDeleteTask_test.json")
	defer cleanup()

	repo := NewTaskRepository(store)

	repo.AppendTaskToDatabase(models.Task{ID: 10, Description: "Task para deletar", IsDone: false})

	err := repo.DeleteTaskFromDatabase(10)

	if err != nil {
		t.Fatalf("Unexpected error while deleting task: %v", err)
	}

	_, err = repo.GetOneTaskFromDatabase(10)

	if err == nil {
		t.Errorf("Expected an error while searching already deleted task, but task it still in memory")
	}
}

func TestMarkTaskAsDone(t *testing.T) {
	cleanup, store := setupTestEnvironment("db_TestMarkTaskAsDone_test.json")
	defer cleanup()

	repo := NewTaskRepository(store)

	repo.AppendTaskToDatabase(models.Task{ID: 1, Description: "Task para marcar como completada", IsDone: false})

	err := repo.MarkTaskAsDone(1)

	if err != nil {
		t.Fatalf("Unexpected error while marking task as done: %v", err)
	}

	taskInDb, _ := repo.GetOneTaskFromDatabase(1)

	if taskInDb.IsDone != true {
		t.Errorf("Task was not set as done correctly")
	}
}

func TestGetLastTaskId(t *testing.T) {
	cleanup, store := setupTestEnvironment("db_TestGetLastTaskId_test.json")
	defer cleanup()

	repo := NewTaskRepository(store)

	repo.AppendTaskToDatabase(models.Task{})

	lastTaskId := repo.GetLastTaskId()

	if lastTaskId != 0 {
		t.Errorf("Expected the last task id to be exactly 0 when there is no tasks added yet.")
	}

	repo.AppendTaskToDatabase(models.Task{ID: 1, Description: "Task 1", IsDone: false})
	repo.AppendTaskToDatabase(models.Task{ID: 1, Description: "Task 2", IsDone: false})
	// ...
	repo.AppendTaskToDatabase(models.Task{ID: 100, Description: "Task 100", IsDone: false})

	lastTaskId = repo.GetLastTaskId()

	if lastTaskId != 100 {
		t.Errorf("Expected the last task id to be exactly 100")
	}
}

func TestLoadDatabaseToMemory(t *testing.T) {
	cleanup, store := setupTestEnvironment("db_TestGetLastTaskId_test.json")
	defer cleanup()

	repo := NewTaskRepository(store)

	repo.store.WriteDatabase([]models.Task{
		{ID: 1, Description: "Task 1", IsDone: false},
	})

	err := repo.LoadDatabaseToMemory()

	if err != nil {
		t.Fatalf("Unexpected error while loading the database file to the memory: %v", err)
	}

	if len(repo.cachedTasks) != 1 {
		t.Errorf("Expected exactly 1 task loaded in memory, found %d", len(repo.cachedTasks))
	}
}
