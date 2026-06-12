package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func startHttpApi() {
	fmt.Println("Servidor http iniciado")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", httpServerStatus)
	mux.HandleFunc("POST /task", httpCreateTask)
	mux.HandleFunc("GET /task", httpGetAllTasks)
	mux.HandleFunc("PATCH /task/{id}/done", httpMarkTaskAsDone)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", mux)
}

func httpServerStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func httpCreateTask(w http.ResponseWriter, r *http.Request) {
	var task CreateTaskDto

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Failed to decode request", http.StatusUnprocessableEntity)
		return
	}

	lastTaskId, err := getLastTaskId()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskId := lastTaskId + 1

	appendTaskToDatabase(Task{taskId, task.Description, false})

	w.WriteHeader(http.StatusCreated)
}

func httpMarkTaskAsDone(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(w, "Failed to convert string to numer", http.StatusUnprocessableEntity)
		return
	}

	err = markTaskAsDone(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func httpGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := getAllTasksFromDatabase()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
