package main

import (
	"fmt"
	"log"
	"os"

	"chrisldo.com/todo-cli/internal/todo/api"
	"chrisldo.com/todo-cli/internal/todo/cli"
	"chrisldo.com/todo-cli/internal/todo/db"
	"chrisldo.com/todo-cli/internal/todo/repository"
)

func main() {
	databaseUrl := "db.json"
	dataBase, err := db.NewFileStore(databaseUrl)

	if err != nil {
		panic(err)
	}

	repo := repository.NewTaskRepository(dataBase)

	err = repo.LoadDatabaseToMemory()

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "api" {
		httpServer := api.NewRestApi(repo)
		log.Fatal(httpServer.StartHttpApi())
	}

	if len(os.Args) == 1 {
		interactiveCli := cli.NewTaskInteractiveCli(repo)
		interactiveCli.StartInteractiveMenu()
		return
	}

	taskCli := cli.NewTaskCli(repo)
	taskCli.HandleCli(os.Args)
}
