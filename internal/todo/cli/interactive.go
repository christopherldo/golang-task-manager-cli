package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"chrisldo.com/todo-cli/internal/todo/models"
	"chrisldo.com/todo-cli/internal/todo/repository"
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

func StartInteractiveMenu() {
	scanner := bufio.NewScanner(os.Stdin)
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

			err := repository.AppendTaskToDatabase(models.Task{ID: repository.GetLastTaskId() + 1, Description: taskDescription, IsDone: false})

			if err != nil {
				fmt.Printf("Error trying to append task to the database: %s", err.Error())
				continue
			}

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
			taskList := repository.GetAllTasksFromDatabase()

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

				err = repository.MarkTaskAsDone(taskToBeMarkedAsCompleted)

				if err != nil {
					fmt.Println("===============================================")
					fmt.Println(err)
					fmt.Println("===============================================")
					continue
				}

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
			fmt.Printf("Erro ao ler a entrada: %s", err.Error())
			continue
		}
	}
}
