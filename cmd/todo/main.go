package main

import (
	"fmt"
	"os"

	"chrisldo.com/todo-cli/internal/todo/cli"
	"chrisldo.com/todo-cli/internal/todo/db"
	"chrisldo.com/todo-cli/internal/todo/repository"
)

func main() {
	err := db.EnsureDataBaseExists()

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	err = repository.LoadDatabaseToMemory()

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	allArgs := os.Args

	if len(allArgs) == 1 {
		cli.StartInteractiveMenu()
		return
	}

	cli.HandleCli(allArgs)
}
