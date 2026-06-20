package cli

import (
	"fmt"
	"strconv"

	"chrisldo.com/todo-cli/internal/todo/models"
	"chrisldo.com/todo-cli/internal/todo/repository"
)

type TaskCli struct {
	repo *repository.TaskRepository
}

func NewTaskCli(repo *repository.TaskRepository) *TaskCli {
	return &TaskCli{
		repo: repo,
	}
}

func (taskCli *TaskCli) HandleCli(allArgs []string) {
	switch allArgs[1] {
	case "help":
		taskCli.cliFuncHelp()
	case "add":
		taskCli.cliFuncAdd(allArgs)
	case "list":
		taskCli.cliFuncList()
	case "done":
		taskCli.cliFuncDone(allArgs)
	case "edit":
		taskCli.cliFuncEdit(allArgs)
	case "delete":
		taskCli.cliFuncDelete(allArgs)
	default:
		fmt.Println("Opção inválida")
	}
}

func (taskCli *TaskCli) cliFuncHelp() {
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

func (taskCli *TaskCli) cliFuncAdd(args []string) {
	if len(args) < 3 {
		fmt.Println("Erro: Você precisa fornecer a descrição da task. Ex: todo-cli add \"Minha task\"")
		return
	}

	taskDescription := args[2]

	err := taskCli.repo.AppendTaskToDatabase(models.Task{ID: taskCli.repo.GetLastTaskId() + 1, Description: taskDescription, IsDone: false})

	if err != nil {
		fmt.Printf("Error trying to append task to the databse: %s", err.Error())
		return
	}

	fmt.Println(`===============================================
Task adicionada!
===============================================`)
}

func (taskCli *TaskCli) cliFuncList() {
	fmt.Println("===============================================")

	taskList := taskCli.repo.GetAllTasksFromDatabase()

	if len(taskList) == 0 {
		fmt.Println("Nenhuma task ainda adicionada")
	} else {
		for _, task := range taskList {
			fmt.Printf("ID: %d, Descrição: %s, Completada? %t\n", task.ID, task.Description, task.IsDone)
		}
	}

	fmt.Println("===============================================")
}

func (taskCli *TaskCli) cliFuncDone(args []string) {
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

	err = taskCli.repo.MarkTaskAsDone(taskId)

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

func (taskCli *TaskCli) cliFuncEdit(args []string) {
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

	task, err := taskCli.repo.GetOneTaskFromDatabase(taskId)

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

	err = taskCli.repo.UpdateTaskOnDatabase(task)

	if err != nil {
		fmt.Printf("Failed to update the task on the database.")
		return
	}
}

func (taskCli *TaskCli) cliFuncDelete(args []string) {
	if len(args) < 3 {
		fmt.Println("Erro: Você precisa fornecer o ID da task. Ex: todo-cli delete 1")
		return
	}

	taskIdString := args[2]
	taskId, err := strconv.Atoi(taskIdString)

	if err != nil {
		fmt.Printf("Failed to parse %s to int", taskIdString)
		return
	}

	err = taskCli.repo.DeleteTaskFromDatabase(taskId)

	if err != nil {
		fmt.Println("===============================================")
		fmt.Println(err)
		fmt.Println("===============================================")
		return
	}

	fmt.Println("===============================================")
	fmt.Printf("Task de ID: %d deletada com sucesso\n", taskId)
	fmt.Println("===============================================")
}
