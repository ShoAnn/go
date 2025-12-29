package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// define resource struct
type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var (
	todolist = []Task{
		{ID: "1", Title: "Go to the grocery store", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", Title: "Go to the capital", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	currentId    = 1
	todolistLock sync.Mutex
)

func main() {
	// setup env
	// setup db conn
	// define mux
	// define routes
	// start server
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database ping successful!")
	defer dbpool.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", getAllTasks)
	mux.HandleFunc("GET /tasks/{id}", getTask)
	mux.HandleFunc("POST /tasks", createTask)
	mux.HandleFunc("PUT /tasks/{id}", updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", deleteTask)

	fmt.Println("Starting server...")
	if err := http.ListenAndServe("localhost:8090", mux); err != nil {
		fmt.Println(err.Error())
	}
}

// handler functions
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	// val (if empty return msg to w)
	// write header
	// encode() all resource data
	if len(todolist) != 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todolist)
	} else {
		fmt.Fprintln(w, "your todolist is currently empty")
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	// get id from url
	// val (if id empty return msg to w)
	// declare new instance of the resource struct for the returned data
	// search task with id (idiom use var found)
	// if empty return
	// write header and encode found data WEH

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "id cannot empty", http.StatusBadRequest)
		return
	}

	var taskFound Task

	var found bool
	for _, task := range todolist {
		if task.ID == idStr {
			taskFound = task
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&taskFound); err != nil {
		http.Error(w, "error encoding to json", http.StatusInternalServerError)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	// create new instance of the resource struct
	// request validation
	// increment id
	// lock > append > unlock
	// write header
	// write header and encode found data WEH
	var newTask Task

	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
	}
	if newTask.ID != "" {
		http.Error(w, "ID field must not be set by the client", http.StatusBadRequest)
		return
	}

	newTask.ID = strconv.Itoa(currentId)
	currentId++

	todolistLock.Lock()
	defer todolistLock.Unlock()
	todolist = append(todolist, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		http.Error(w, "error encoding to json", http.StatusInternalServerError)
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	// get id from url
	// val (if id empty return msg to w)
	// declare new instance of the resource struct for the edited data
	// search data with the id (idiom use var found)
	// if empty return
	// update data (attention to unchanged fields)
	// write header and encode found data WEH

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "id cannot empty", http.StatusBadRequest)
		return
	}

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "error decoding request body ", http.StatusBadRequest)
	}

	todolistLock.Lock()
	defer todolistLock.Unlock()

	foundId := -1
	for i, task := range todolist {
		if task.ID == idStr {
			foundId = i
			break
		}
	}

	if foundId == -1 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	todolist[foundId].Title = newTask.Title
	todolist[foundId].Completed = newTask.Completed

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		http.Error(w, "error encoding to json", http.StatusInternalServerError)
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	// get id from url
	// val (if id empty return msg to w)
	// search data with the id (idiom use var found)
	// if not found return error
	// delete
	// write header

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "id cannot empty", http.StatusBadRequest)
		return
	}

	todolistLock.Lock()
	defer todolistLock.Unlock()

	foundId := -1
	for i, task := range todolist {
		if task.ID == idStr {
			foundId = i
			break
		}
	}

	if foundId == -1 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	todolist = append(todolist[:foundId], todolist[foundId+1:]...)

	w.WriteHeader(http.StatusNoContent)
}
