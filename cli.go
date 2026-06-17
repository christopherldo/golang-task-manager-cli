package main

import (
	"fmt"
	"strconv"
)

func cliFuncHelp() {
	fmt.Println(`===============================================
TASK CLI - Guia de Uso
===============================================
Você pode usar este programa de duas formas:

1. MODO INTERATIVO:
Basta rodar o programa sem nenhum argumento.
Exemplo: ./todo-cli

2. MODO COMANDOS DIRETOS (CLI):
Passe a ação e os argumentos diretamente no terminal.
Uso: ./todo-cli [comando] [argumento]

Comandos disponíveis:
  add <descrição>   Adiciona uma nova task.
                    Ex: ./todo-cli add "Estudar Go"

  list              Lista todas as tasks salvas.
                    Ex: ./todo-cli list

  done <id>         Marca a task correspondente como concluída.
                    Ex: ./todo-cli done 1

  help              Exibe este menu de ajuda.

  api               Inicia um servidor http para rodar o programa.
===============================================`)
}

func cliFuncAdd(args []string) {
	if len(args) < 3 {
		fmt.Println("Erro: Você precisa fornecer a descrição da task. Ex: todo-cli add \"Minha task\"")
		return
	}

	taskDescription := args[2]
	lastTaskId := getLastTaskId()

	taskId := lastTaskId + 1

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

func cliFuncApi() {
	startHttpApi()
}
