package cli

import (
	"fmt"
	"strconv"

	"chrisldo.com/todo-cli/internal/todo/api"
	"chrisldo.com/todo-cli/internal/todo/models"
	"chrisldo.com/todo-cli/internal/todo/repository"
)

func HandleCli(allArgs []string) {
	switch allArgs[1] {
	case "help":
		cliFuncHelp()
	case "add":
		cliFuncAdd(allArgs)
	case "list":
		cliFuncList()
	case "done":
		cliFuncDone(allArgs)
	case "edit":
		cliFuncEdit(allArgs)
	case "api":
		cliFuncApi()
	default:
		fmt.Println("Opção inválida")
	}
}

func cliFuncHelp() {
	fmt.Println(`=======================================================================================
TASK CLI - Guia de Uso
=======================================================================================
Você pode usar este programa de duas formas:

1. MODO INTERATIVO:
Basta rodar o programa sem nenhum argumento.
Exemplo: ./todo-cli

2. MODO COMANDOS DIRETOS (CLI):
Passe a ação e os argumentos diretamente no terminal.
Uso: ./todo-cli [comando] [argumento]

Comandos disponíveis:
  add <descrição>                       Adiciona uma nova task.
                                        Ex: ./todo-cli add "Estudar Go"

  list                                  Lista todas as tasks salvas.
                                        Ex: ./todo-cli list

  update <id> <description> <isDone?>   Atualiza uma task.
                                        Ex: ./todo-cli update 1 "Estudar Go"
                                        Ex: ./todo-cli update 1 "Estudar Go" true

  done <id>                             Marca a task correspondente como concluída.
                                        Ex: ./todo-cli done 1

  help                                  Exibe este menu de ajuda.

  api                                   Inicia um servidor http para rodar o programa.
=======================================================================================`)
}

func cliFuncAdd(args []string) {
	if len(args) < 3 {
		fmt.Println("Erro: Você precisa fornecer a descrição da task. Ex: todo-cli add \"Minha task\"")
		return
	}

	taskDescription := args[2]

	err := repository.AppendTaskToDatabase(models.Task{ID: repository.GetLastTaskId() + 1, Description: taskDescription, IsDone: false})

	if err != nil {
		fmt.Printf("Error trying to append task to the databse: %s", err.Error())
		return
	}

	fmt.Println(`===============================================
Task adicionada!
===============================================`)
}

func cliFuncList() {
	fmt.Println("===============================================")

	taskList := repository.GetAllTasksFromDatabase()

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
	if len(args) < 3 {
		fmt.Println("Erro: Você precisa fornecer o ID da task. Ex: todo-cli done 1")
		return
	}

	taskIdString := args[2]
	taskId, err := strconv.Atoi(taskIdString)

	if err != nil {
		fmt.Printf("Failed to parse %s to int", taskIdString)
		return
	}

	err = repository.MarkTaskAsDone(taskId)

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

func cliFuncEdit(args []string) {
	if len(args) < 4 {
		fmt.Println("Erro: Você precisa fornecer o ID da task e a descrição. Ex: todo-cli edit 1 \"Estudar Go\"")
		return
	}

	taskIdString := args[2]
	taskDescriptionString := args[3]

	var taskIsDoneString string

	if len(args) >= 5 {
		taskIsDoneString = args[4]
	}

	taskId, err := strconv.Atoi(taskIdString)

	if err != nil {
		fmt.Printf("Failed to parse %s to int", taskIdString)
		return
	}

	task, err := repository.GetOneTaskFromDatabase(taskId)

	if err != nil {
		fmt.Printf("Failed to find task with ID: %d", taskId)
		return
	}

	task.Description = taskDescriptionString

	if taskIsDoneString != "" {
		var newTaskIsDoneStatusBool bool

		switch taskIsDoneString {
		case "true":
			newTaskIsDoneStatusBool = true
		case "false":
			newTaskIsDoneStatusBool = false
		default:
			fmt.Printf("Invalid option for task status")
			return
		}

		task.IsDone = newTaskIsDoneStatusBool
	}

	err = repository.UpdateTaskOnDatabase(task)

	if err != nil {
		fmt.Printf("Failed to update the task on the database.")
		return
	}
}

func cliFuncApi() {
	err := api.StartHttpApi()

	if err != nil {
		fmt.Printf("Error on cli function that starts api: %s", err.Error())
		return
	}
}
