package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"chrisldo.com/todo-cli/internal/todo/dto"
	"chrisldo.com/todo-cli/internal/todo/models"
	"chrisldo.com/todo-cli/internal/todo/repository"
)

func StartHttpApi() error {
	fmt.Println("Servidor http iniciado")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", httpServerStatus)
	mux.HandleFunc("POST /task", httpCreateTask)
	mux.HandleFunc("GET /task", httpGetAllTasks)
	mux.HandleFunc("PATCH /task/{id}/done", httpMarkTaskAsDone)
	mux.HandleFunc("PUT /task/{id}", httpEditTask)

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		return fmt.Errorf("Error starting HTTP API: %w", err)
	}

	return nil
}

func httpServerStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func httpCreateTask(w http.ResponseWriter, r *http.Request) {
	var task dto.CreateTaskDto

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Failed to decode request", http.StatusUnprocessableEntity)
		return
	}

	err := repository.AppendTaskToDatabase(models.Task{ID: repository.GetLastTaskId() + 1, Description: task.Description, IsDone: false})

	if err != nil {
		http.Error(w, "Internal Server Error while trying to append task to the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func httpMarkTaskAsDone(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(w, "Failed to convert string to numer", http.StatusUnprocessableEntity)
		return
	}

	err = repository.MarkTaskAsDone(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func httpGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := repository.GetAllTasksFromDatabase()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(tasks)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func httpEditTask(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(w, "Failed to convert string to numer", http.StatusUnprocessableEntity)
		return
	}

	var taskToBeUpdated dto.EditTaskDto

	if err := json.NewDecoder(r.Body).Decode(&taskToBeUpdated); err != nil {
		http.Error(w, "Failed to decode request", http.StatusUnprocessableEntity)
		return
	}

	task, err := repository.GetOneTaskFromDatabase(id)

	if err != nil {
		http.Error(w, "Task not found with this ID", http.StatusNotFound)
		return
	}

	task.Description = taskToBeUpdated.Description

	if taskToBeUpdated.IsDone != nil {
		task.IsDone = *taskToBeUpdated.IsDone
	}

	err = repository.UpdateTaskOnDatabase(task)

	if err != nil {
		http.Error(w, "Internal Server Error while trying to update task on the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
