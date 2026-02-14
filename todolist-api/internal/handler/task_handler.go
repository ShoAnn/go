package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ShoAnn/go-playground/todolist-api/internal/domain"
	"github.com/go-playground/validator/v10"
)

type TaskHandler struct {
	service  domain.TaskService
	validate *validator.Validate
}

func NewTaskHandler(s domain.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	taskList, err := h.service.ListTasks(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(taskList)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
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
		return
	}
	// We pass r.Context() so if the user cancels the request, the DB stops working too.
	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		// If pgx can't find the row, it returns a specific "No Rows" error
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var params domain.CreateTaskParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(params); err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(r.Context(), &domain.CreateTaskParams{
		Title: params.Title,
	})
	if err != nil {
		http.Error(w, "error creating task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	var params domain.UpdateTaskParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid request body ", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(params); err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	updatedTask, err := h.service.EditTask(r.Context(), id, &domain.UpdateTaskParams{
		Title:     params.Title,
		Completed: params.Completed,
		Version:   params.Version,
	})
	if err != nil {
		http.Error(w, "error updating task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updatedTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.DeleteTask(r.Context(), id)
	if err != nil {
		http.Error(w, "error deleting task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
