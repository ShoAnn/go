package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	defer dbpool.Close()
	if err := dbpool.Ping(ctx); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	// 4. Wrap the pool with sqlc
	// database.New comes from the db.go file sqlc generated for you
	queries := db.New(dbpool)
	server := &Server{
		Queries: queries,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", server.getAllTasks)
	mux.HandleFunc("GET /tasks/{id}", server.getTask)
	mux.HandleFunc("POST /tasks", server.createTask)
	mux.HandleFunc("PUT /tasks/{id}", server.updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", server.deleteTask)

	fmt.Println("Starting server...")
	if err := http.ListenAndServe("localhost:8090", mux); err != nil {
		fmt.Println(err.Error())
	}
}
