package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"sync"
)

var cacheMutex sync.RWMutex
var cachedTasks = []Task{}
var cachedLastTaskId *int

func appendTaskToDatabase(taskToBeAdded Task) error {
	cacheMutex.Lock()

	cachedTasks = append(cachedTasks, taskToBeAdded)
	tasksToWrite := cachedTasks
	cachedLastTaskId = &taskToBeAdded.ID

	cacheMutex.Unlock()

	err := writeDatabase(tasksToWrite)

	if err != nil {
		return fmt.Errorf("Error appeding task to the database: %w", err)
	}

	return nil
}

func getAllTasksFromDatabase() []Task {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	return cachedTasks
}

/*
func updateTaskOnDatabase(database *os.File, taskToBeUpdated Task) {
	allTasks := getAllTasksFromDatabase()

	idx := slices.IndexFunc(allTasks, func(task Task) bool {
		return task.ID == taskToBeUpdated.ID
	})

	if idx == -1 {
		fmt.Println("Task não encontrado com esse ID")
		return
	}

	allTasks[idx] = taskToBeUpdated
	writeDatabase(database, allTasks)
}
*/

func markTaskAsDone(taskId int) error {
	cacheMutex.Lock()

	idx := slices.IndexFunc(cachedTasks, func(task Task) bool {
		return task.ID == taskId
	})

	if idx == -1 {
		cacheMutex.Unlock()
		return fmt.Errorf("Task não encontrado com esse ID")
	}

	cachedTasks[idx].IsDone = true
	tasksToWrite := cachedTasks

	cacheMutex.Unlock()

	err := writeDatabase(tasksToWrite)

	if err != nil {
		return fmt.Errorf("Error marking task as done: %w", err)
	}

	return nil
}

func getLastTaskId() int {
	if cachedLastTaskId != nil {
		cacheMutex.RLock()
		defer cacheMutex.RUnlock()

		return *cachedLastTaskId
	}

	tasks := getAllTasksFromDatabase()

	if len(tasks) == 0 {
		return 0
	}

	lastIndex := len(tasks) - 1

	lastTask := tasks[lastIndex]

	return lastTask.ID
}

func loadDatabaseToMemory() error {
	bytes, err := readDatabase()

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &cachedTasks)

	if err != nil {
		return fmt.Errorf("Error parsing JSON: %w", err)
	}

	lastTaskId := getLastTaskId()
	cachedLastTaskId = &lastTaskId

	return nil
}
