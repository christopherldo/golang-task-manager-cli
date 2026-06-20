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
	EditingTaskSession    ProgramSession = "editing_task"
	DeletingTaskSession   ProgramSession = "deleting_task"
	ExitSession           ProgramSession = "exit"
)

type TaskInteractiveCli struct {
	repo *repository.TaskRepository
}

func NewTaskInteractiveCli(repo *repository.TaskRepository) *TaskInteractiveCli {
	return &TaskInteractiveCli{
		repo: repo,
	}
}

func (interactiveCli *TaskInteractiveCli) StartInteractiveMenu() {
	scanner := bufio.NewScanner(os.Stdin)
	programSession := MenuSession

	for {
		// Exiting the progam
		if programSession == ExitSession {
			fmt.Println(`===============================================
Muito obrigado por usar o TASK CLI.
===============================================`)
			break
		}

		// Menu
		if programSession == MenuSession {
			fmt.Println(`===============================================
Bem-vindo ao TASK CLI`)
		}

		// Adding Task
		if programSession == AddingTaskSession {
			fmt.Println("Digite a descrição da sua task:")
			scanner.Scan()
			taskDescription := scanner.Text()

			err := interactiveCli.repo.AppendTaskToDatabase(models.Task{ID: interactiveCli.repo.GetLastTaskId() + 1, Description: taskDescription, IsDone: false})

			if err != nil {
				fmt.Printf("Error trying to append task to the database: %s", err.Error())
				continue
			}

			programSession = AddedTaskSession
		}

		// Task Added
		if programSession == AddedTaskSession {
			fmt.Println(`===============================================
Task adicionada!`)
		}

		// List all tasks
		if programSession == ListAllTasksSession {
			taskList := interactiveCli.repo.GetAllTasksFromDatabase()

			fmt.Println("===============================================")

			if len(taskList) == 0 {
				fmt.Println("Nenhuma task ainda adicionada")
			} else {
				for _, task := range taskList {
					fmt.Printf("ID: %d, Descrição: %s, Completada? %t\n", task.ID, task.Description, task.IsDone)
				}
			}
		}

		// Setting Task as Completed
		if programSession == CompletingTaskSession {
			for {
				selectedTask, err := selectTaskById(scanner)

				if err != nil {
					fmt.Println("Opção inválida")
					continue
				}

				if selectedTask == 0 {
					programSession = MenuSession
					fmt.Println(`===============================================
Opção cancelada.`)
					break
				}

				err = interactiveCli.repo.MarkTaskAsDone(selectedTask)

				if err != nil {
					fmt.Println("===============================================")
					fmt.Println(err)
					fmt.Println("===============================================")
					continue
				}

				fmt.Println("===============================================")
				fmt.Printf("Task de ID: %d marcada como concluída\n", selectedTask)

				break
			}
		}

		if programSession == EditingTaskSession {
			for {
				selectedTask, err := selectTaskById(scanner)

				if err != nil {
					fmt.Println("Opção inválida")
					continue
				}

				if selectedTask == 0 {
					programSession = MenuSession
					fmt.Println(`===============================================
Opção cancelada.`)
					break
				}

				var taskToBeUpdated models.Task

				taskToBeUpdated, err = interactiveCli.repo.GetOneTaskFromDatabase(selectedTask)

				if err != nil {
					fmt.Println("Task não encontrada. Digite o ID novamente")
					continue
				}

				fmt.Println("Digite a nova descrição da task ou ENTER para manter a antiga:")

				scanner.Scan()
				newTaskDescription := scanner.Text()

				if newTaskDescription == "" {
					newTaskDescription = taskToBeUpdated.Description
				}

				fmt.Println("Task já concluída? (Y/N) ou ENTER para manter o estado antigo:")

				scanner.Scan()
				newTaskIsDoneStatusText := scanner.Text()

				var newTaskIsDoneStatusBool bool

				switch newTaskIsDoneStatusText {
				case "Y":
					newTaskIsDoneStatusBool = true
				case "N":
					newTaskIsDoneStatusBool = false
				case "":
					newTaskIsDoneStatusBool = taskToBeUpdated.IsDone
				}

				taskToBeUpdated.Description = newTaskDescription
				taskToBeUpdated.IsDone = newTaskIsDoneStatusBool

				err = interactiveCli.repo.UpdateTaskOnDatabase(taskToBeUpdated)

				if err != nil {
					fmt.Println("===============================================")
					fmt.Println(err)
					fmt.Println("===============================================")
					continue
				}

				fmt.Println("===============================================")
				fmt.Printf("Task de ID: %d editada com sucesso\n", selectedTask)

				break
			}
		}

		if programSession == DeletingTaskSession {
			for {
				selectedTask, err := selectTaskById(scanner)

				if err != nil {
					fmt.Println("Opção inválida")
					continue
				}

				if selectedTask == 0 {
					programSession = MenuSession
					fmt.Println(`===============================================
Opção cancelada.`)
					break
				}

				err = interactiveCli.repo.DeleteTaskFromDatabase(selectedTask)

				if err != nil {
					fmt.Println("===============================================")
					fmt.Println(err)
					fmt.Println("===============================================")
					continue
				}

				fmt.Println("===============================================")
				fmt.Printf("Task de ID: %d deletada com sucesso\n", selectedTask)

				break
			}
		}

		printMenuCli(programSession)

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
		case 4:
			programSession = EditingTaskSession
		case 5:
			programSession = DeletingTaskSession
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

func selectTaskById(scanner *bufio.Scanner) (int, error) {
	fmt.Println("Digite o ID da task ou 0 para cancelar:")
	scanner.Scan()
	taskToBeSelected, err := strconv.Atoi(scanner.Text())

	if err != nil {
		return 0, fmt.Errorf("Opção inválida")
	}

	if taskToBeSelected == 0 {
		return 0, nil
	}

	return taskToBeSelected, nil
}
