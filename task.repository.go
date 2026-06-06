package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
)

func appendTaskToDatabase(taskToBeAdded Task) {
	allTasks := getAllTasksFromDatabase()

	allTasks = append(allTasks, taskToBeAdded)

	writeDatabase(allTasks)
}

func getAllTasksFromDatabase() []Task {
	bytes := readDatabase()
	var tasks []Task

	err := json.Unmarshal(bytes, &tasks)

	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return tasks
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
	allTasks := getAllTasksFromDatabase()

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

func getLastTaskId() int {
	tasks := getAllTasksFromDatabase()

	if len(tasks) == 0 {
		return 0
	}

	lastIndex := len(tasks) - 1

	lastTask := tasks[lastIndex]

	return lastTask.ID
}
