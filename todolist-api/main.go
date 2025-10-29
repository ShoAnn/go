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

	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Get all tasks")
		// get all tasks code
	})

	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "Get a task with the id %s", id)
		// get 1 task code
	})

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Post a task")
		// create task code
	})

	mux.HandleFunc("PUT /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "Update task with the id %s", id)
		// update code
	})

	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "Delete task with the id %s", id)
		// delete code
	})

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err.Error())
	}
}
