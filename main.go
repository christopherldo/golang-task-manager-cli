package main

import (
	"fmt"
	"os"
)

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
	case "api":
		cliFuncApi()
	default:
		fmt.Println("Opção inválida")
	}
}
