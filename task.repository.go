package main

import (
	"encoding/json"
	"fmt"
	"slices"
)

func appendTaskToDatabase(taskToBeAdded Task) error {
	allTasks, err := getAllTasksFromDatabase()

	if err != nil {
		return err
	}

	allTasks = append(allTasks, taskToBeAdded)

	writeDatabase(allTasks)

	return nil
}

func getAllTasksFromDatabase() ([]Task, error) {
	bytes, err := readDatabase()

	if err != nil {
		return nil, err
	}

	var tasks []Task

	err = json.Unmarshal(bytes, &tasks)

	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %w", err)
	}

	return tasks, err
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
	allTasks, err := getAllTasksFromDatabase()

	if err != nil {
		return err
	}

	idx := slices.IndexFunc(allTasks, func(task Task) bool {
		return task.ID == taskId
	})

	if idx == -1 {
		return fmt.Errorf("Task não encontrado com esse ID")
	}

	allTasks[idx].IsDone = true
	writeDatabase(allTasks)

	return nil
}

func getLastTaskId() (int, error) {
	tasks, err := getAllTasksFromDatabase()

	if err != nil {
		return 0, err
	}

	if len(tasks) == 0 {
		return 0, nil
	}

	lastIndex := len(tasks) - 1

	lastTask := tasks[lastIndex]

	return lastTask.ID, err
}
