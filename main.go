package main

import (
	"fmt"
	"os"
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

const (
	DB_URL = "db.json"
)

type Task struct {
	ID          int
	Description string
	IsDone      bool
}

func main() {
	ensureDataBaseExists()

	allArgs := os.Args

	if len(allArgs) == 1 {
		startInteractiveMenu()
		return
	}

	switch allArgs[1] {
	case "help":
		cliFuncHelp()
	case "add":
		cliFuncAdd(allArgs)
	case "list":
		cliFuncList()
	case "done":
		cliFuncDone(allArgs)
	default:
		fmt.Println("Opção inválida")
	}
}
