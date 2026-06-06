package main

import (
	"fmt"
	"log"
	"strconv"
)

func cliFuncHelp() {
	fmt.Println("Print menu")
}

func cliFuncAdd(args []string) {
	taskDescription := args[2]
	taskId := getLastTaskId() + 1

	appendTaskToDatabase(Task{taskId, taskDescription, false})
	fmt.Println(`===============================================
Task adicionada!
===============================================`)
}

func cliFuncList() {
	fmt.Println("===============================================")

	taskList := getAllTasksFromDatabase()

	if len(taskList) == 0 {
		fmt.Println("Nenhuma task ainda adicionada")
	} else {
		for _, task := range taskList {
			fmt.Printf("ID: %d, Descrição: %s, Completada? %t\n", task.ID, task.Description, task.IsDone)
		}
	}

	fmt.Println("===============================================")
}

func cliFuncDone(args []string) {
	taskIdString := args[2]
	taskId, err := strconv.Atoi(taskIdString)

	if err != nil {
		log.Fatalf("Failed to parse %s to int", taskIdString)
	}

	err = markTaskAsDone(taskId)

	if err != nil {
		fmt.Println("===============================================")
		fmt.Println(err)
		fmt.Println("===============================================")
		return
	}

	fmt.Println("===============================================")
	fmt.Printf("Task de ID: %d marcada como concluída\n", taskId)
	fmt.Println("===============================================")
}
