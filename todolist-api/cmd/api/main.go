package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ShoAnn/go-playground/todolist-api/internal/handler"
	"github.com/ShoAnn/go-playground/todolist-api/internal/repository"
	"github.com/ShoAnn/go-playground/todolist-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// setup env
	// setup db conn
	// define mux
	// define routes
	// start server
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()
	if err := dbpool.Ping(ctx); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	taskRepo := repository.NewTaskRepository(dbpool)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", taskHandler.ListTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTask)
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("PATCH /tasks/{id}", taskHandler.CompleteTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	log.Println("Starting server...")
	if err := http.ListenAndServe("localhost:8090", mux); err != nil {
		fmt.Println(err.Error())
	}
}
