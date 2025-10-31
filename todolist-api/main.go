package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	type Task struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	mux.HandleFunc("GET /tasks", getTasks)
	mux.HandleFunc("GET /tasks/{id}", getTask)
	mux.HandleFunc("POST /tasks", createTask)
	mux.HandleFunc("PUT /tasks/{id}", updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", deleteTask)

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err.Error())
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get all tasks")
	// get all tasks code
}
func getTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Get a task with the id %s", id)
	// get 1 task code
}
func createTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Post a task")
	// create task code
}
func updateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Update task with the id %s", id)
	// update code
}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Delete task with the id %s", id)
	// delete code
}
