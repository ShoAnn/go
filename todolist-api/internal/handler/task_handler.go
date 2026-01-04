package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func (s *Server) getAllTasks(w http.ResponseWriter, r *http.Request) {
	// val (if empty return msg to w)
	// write header
	// encode() all resource data

	taskList, err := s.Queries.GetAllTasks(r.Context())
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(taskList); err != nil {
		log.Printf("error encoding to json: %v", err)
	}
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) {
	// get id from url
	// val (if id empty return msg to w)
	// declare new instance of the resource struct for the returned data
	// search task with id (idiom use var found)
	// if empty return
	// write header and encode found data WEH

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}
	// sqlc generated this method!
	// We pass r.Context() so if the user cancels the request, the DB stops working too.
	task, err := s.Queries.GetTask(r.Context(), int32(id))
	if err != nil {
		// If pgx can't find the row, it returns a specific "No Rows" error
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("error encoding to json: %v", err)
	}
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	// create new instance of the resource struct
	// request validation
	// increment id
	// lock > append > unlock
	// write header
	// write header and encode found data WEH
	type request struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
	}

	task, err := s.Queries.CreateTask(r.Context(), db.CreateTaskParams{
		Title:     req.Title,
		Completed: req.Completed,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "error creating task", http.StatusNotFound)
			return
		}
		http.Error(w, "error creating task", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(task); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) {
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
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	type request struct {
		Id        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body ", http.StatusBadRequest)
		return
	}

	updateTask, err := s.Queries.UpdateTask(r.Context(), db.UpdateTaskParams{
		ID:        int32(id),
		Title:     req.Title,
		Completed: req.Completed,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "error creating task", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updateTask); err != nil {
		log.Printf("error encoding to json: %v", err)
	}
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	deleteResult, err := s.Queries.DeleteTask(r.Context(), int32(id))
	rowsAffected := deleteResult.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error deleting task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
