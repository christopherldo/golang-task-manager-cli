package main

import (
	"fmt"
	"os"
)

func main() {
	err := ensureDataBaseExists()

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	err = loadDatabaseToMemory()

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

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
	case "api":
		cliFuncApi()
	default:
		fmt.Println("Opção inválida")
	}
}
