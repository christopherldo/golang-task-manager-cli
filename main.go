// Fase 1: O Básico (Estruturas de Dados)

// Comece criando o programa para rodar apenas na memória enquanto o terminal estiver aberto.

//     O que fazer: Crie um menu simples no terminal onde você possa adicionar uma tarefa, listar todas as tarefas e marcar uma como concluída.

//     O que você vai aprender:

//         Criação e uso de structs (ex: criar uma estrutura Task com ID, Description e IsDone).

//         Trabalhar com Slices (as listas do Go) para armazenar as tarefas.

//         Laços de repetição (for) e condicionais (if/else).

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type ProgramSession string

const (
	MenuSession           ProgramSession = "menu"
	AddingTaskSession     ProgramSession = "adding_task"
	AddedTaskSession      ProgramSession = "added_task"
	ListAllTasksSession   ProgramSession = "list_all_tasks"
	CompletingTaskSession ProgramSession = "completing_task"
	CompletedTaskSession  ProgramSession = "completed_task"
	ExitSession           ProgramSession = "exit"
)

type Task struct {
	ID          int
	Description string
	IsDone      bool
}

func main() {
	var taskList []Task

	scanner := bufio.NewScanner(os.Stdin)
	taskId := 1
	programSession := MenuSession

	for {
		if programSession == ExitSession {
			fmt.Println(`===============================================
Muito obrigado por usar o TASK CLI.
===============================================`)
			break
		}

		if programSession == MenuSession {
			fmt.Println(`===============================================
Bem-vindo ao TASK CLI
===============================================
O que deseja?
1 - Adicionar uma task.
2 - Listar todas as tasks adicionadas.
3 - Marcar uma task como concluída.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
		}

		if programSession == AddingTaskSession {
			fmt.Println("Digite a descrição da sua task:")
			scanner.Scan()
			taskDescription := scanner.Text()
			taskList = append(taskList, Task{taskId, taskDescription, false})
			taskId++
			programSession = AddedTaskSession
		}

		if programSession == AddedTaskSession {
			fmt.Println(`===============================================
Task adicionada!
===============================================
O que deseja?
1 - Adicionar outra task.
2 - Listar todas as tasks adicionadas.
3 - Marcar uma task como concluída.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
		}

		if programSession == ListAllTasksSession {
			fmt.Println("===============================================")

			if len(taskList) == 0 {
				fmt.Println("Nenhuma task ainda adicionada")
			} else {
				for _, task := range taskList {
					fmt.Printf("ID: %d, Descrição: %s, Completada? %t\n", task.ID, task.Description, task.IsDone)
				}
			}

			fmt.Println(`===============================================
O que deseja agora?
1 - Adicionar uma task.
2 - Listar novamente todas as tasks adicionadas.
3 - Marcar uma task como concluída.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
		}

		if programSession == CompletingTaskSession {
			for {
				fmt.Println("Digite o ID da task ou 0 para cancelar:")
				scanner.Scan()
				taskToBeMarkedAsCompleted, err := strconv.Atoi(scanner.Text())

				if err != nil {
					fmt.Println("Opção inválida")
					continue
				}

				if taskToBeMarkedAsCompleted == 0 {
					programSession = MenuSession
					fmt.Println(`===============================================
Opção cancelada.
===============================================
O que deseja agora?
1 - Adicionar uma task.
2 - Listar todas as tasks adicionadas.
3 - Marcar outra task como concluída.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
					break
				}

				idx := slices.IndexFunc(taskList, func(task Task) bool {
					return task.ID == taskToBeMarkedAsCompleted
				})

				if idx == -1 {
					fmt.Println("Task não encontrado com esse ID")
				} else {
					taskList[idx].IsDone = true

					fmt.Println("===============================================")
					fmt.Printf("Task de ID: %d marcada como concluída\n", taskToBeMarkedAsCompleted)
					fmt.Println(`===============================================
O que deseja agora?
1 - Adicionar uma task.
2 - Listar todas as tasks adicionadas.
3 - Marcar outra task como concluída.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
					break
				}
			}
		}

		scanner.Scan()
		menuOption, err := strconv.Atoi(scanner.Text())

		if err != nil {
			fmt.Println(`===============================================
Opção inválida!
===============================================`)
			continue
		}

		switch menuOption {
		case 1:
			programSession = AddingTaskSession
		case 2:
			programSession = ListAllTasksSession
		case 3:
			programSession = CompletingTaskSession
		case 9:
			programSession = MenuSession
		case 0:
			programSession = ExitSession
		default:
			fmt.Println("Opção inválida")
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Erro ao ler a entrada:", err)
		}
	}
}
