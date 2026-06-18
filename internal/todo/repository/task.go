package repository

import (
	"encoding/json"
	"fmt"
	"slices"
	"sync"

	"chrisldo.com/todo-cli/internal/todo/db"
	"chrisldo.com/todo-cli/internal/todo/models"
)

var cacheMutex sync.RWMutex
var cachedTasks = []models.Task{}
var cachedLastTaskId *int

func AppendTaskToDatabase(taskToBeAdded models.Task) error {
	cacheMutex.Lock()

	cachedTasks = append(cachedTasks, taskToBeAdded)
	tasksToWrite := cachedTasks
	cachedLastTaskId = &taskToBeAdded.ID

	cacheMutex.Unlock()

	err := db.WriteDatabase(tasksToWrite)

	if err != nil {
		return fmt.Errorf("Error appeding task to the database: %w", err)
	}

	return nil
}

func GetOneTaskFromDatabase(taskId int) (models.Task, error) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	idx := slices.IndexFunc(cachedTasks, func(task models.Task) bool {
		return task.ID == taskId
	})

	if idx == -1 {
		return models.Task{}, fmt.Errorf("Task não encontrado com esse ID")
	}

	return cachedTasks[idx], nil
}

func GetAllTasksFromDatabase() []models.Task {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	return cachedTasks
}

func UpdateTaskOnDatabase(taskToBeUpdated models.Task) error {
	cacheMutex.Lock()

	idx := slices.IndexFunc(cachedTasks, func(task models.Task) bool {
		return task.ID == taskToBeUpdated.ID
	})

	if idx == -1 {
		cacheMutex.Unlock()
		return fmt.Errorf("Task não encontrado com esse ID")
	}

	cachedTasks[idx] = taskToBeUpdated

	tasksToWrite := cachedTasks

	cacheMutex.Unlock()

	err := db.WriteDatabase(tasksToWrite)

	if err != nil {
		return fmt.Errorf("Error while saving task do the database: %w", err)
	}

	return nil
}

func MarkTaskAsDone(taskId int) error {
	cacheMutex.Lock()

	idx := slices.IndexFunc(cachedTasks, func(task models.Task) bool {
		return task.ID == taskId
	})

	if idx == -1 {
		cacheMutex.Unlock()
		return fmt.Errorf("Task não encontrado com esse ID")
	}

	cachedTasks[idx].IsDone = true
	tasksToWrite := cachedTasks

	cacheMutex.Unlock()

	err := db.WriteDatabase(tasksToWrite)

	if err != nil {
		return fmt.Errorf("Error marking task as done: %w", err)
	}

	return nil
}

func GetLastTaskId() int {
	if cachedLastTaskId != nil {
		cacheMutex.RLock()
		defer cacheMutex.RUnlock()

		return *cachedLastTaskId
	}

	tasks := GetAllTasksFromDatabase()

	if len(tasks) == 0 {
		return 0
	}

	lastIndex := len(tasks) - 1

	lastTask := tasks[lastIndex]

	return lastTask.ID
}

func LoadDatabaseToMemory() error {
	bytes, err := db.ReadDatabase()

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &cachedTasks)

	if err != nil {
		return fmt.Errorf("Error parsing JSON: %w", err)
	}

	lastTaskId := GetLastTaskId()
	cachedLastTaskId = &lastTaskId

	return nil
}
