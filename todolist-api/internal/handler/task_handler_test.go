package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ShoAnn/go-playground/todolist-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func StructToReader(v interface{}) (io.Reader, error) {
	// Encode the struct to a byte slice
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}
	// Return a new io.Reader from the byte slice
	return bytes.NewReader(b), nil
}

func TestListTasks(t *testing.T) {
	tasks := []*domain.Task{
		{
			ID:        1,
			Title:     "Shopping",
			Completed: false,
			CreatedAt: time.Now().Add(-2 * time.Minute),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		{
			ID:        2,
			Title:     "Work",
			Completed: true,
			CreatedAt: time.Now().Add(-40 * time.Minute),
			UpdatedAt: time.Now(),
			Version:   1,
		},
	}
	jsonTasks, _ := json.MarshalIndent(tasks, "", "  ")
	stringTasks := string(jsonTasks)

	tests := []struct {
		name           string
		mockSetup      func(m *MockTaskService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Happy_path",
			mockSetup: func(m *MockTaskService) {
				m.On("ListTasks", mock.Anything).Return(tasks, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   stringTasks,
		},
		{
			name: "empty list",
			mockSetup: func(m *MockTaskService) {
				m.On("ListTasks", mock.Anything).Return(nil, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockTaskService)
			tc.mockSetup(mockService)
			h := &TaskHandler{service: mockService}

			r := httptest.NewRequest("GET", "/tasks", nil)
			w := httptest.NewRecorder()

			h.ListTasks(w, r)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	task := &domain.Task{
		ID:        1,
		Title:     "Shopping",
		Completed: false,
		CreatedAt: time.Now().Add(-2 * time.Minute),
		UpdatedAt: time.Now(),
		Version:   0,
	}
	jsonTask, _ := json.Marshal(task)
	stringTask := string(jsonTask)

	tests := []struct {
		name           string
		taskId         string
		mockSetup      func(m *MockTaskService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Happy_path",
			taskId: "1",
			mockSetup: func(m *MockTaskService) {
				m.On("GetTask", mock.Anything, 1).Return(task, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   stringTask,
		},
		{
			name:   "id not found",
			taskId: "999",
			mockSetup: func(m *MockTaskService) {
				m.On("GetTask", mock.Anything, 999).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockTaskService)
			tc.mockSetup(mockService)
			h := &TaskHandler{service: mockService}

			mux := http.NewServeMux()
			mux.HandleFunc("GET /tasks/{id}", h.GetTask)
			r := httptest.NewRequest("GET", "/tasks/"+tc.taskId, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	createParams := &domain.CreateTaskParams{Title: "Shopping"}
	task := &domain.Task{
		ID:        1,
		Title:     "Shopping",
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   0,
	}

	tests := []struct {
		name           string
		mockSetup      func(m *MockTaskService)
		expectedStatus int
		expectedTitle  string
	}{
		{
			name: "Happy_path",
			mockSetup: func(m *MockTaskService) {
				m.On("CreateTask", mock.Anything, createParams).Return(task, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedTitle:  task.Title,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockTaskService)
			tc.mockSetup(mockService)
			h := &TaskHandler{service: mockService}
			mux := http.NewServeMux()
			mux.HandleFunc("POST /tasks", h.CreateTask)
			reqBody, _ := StructToReader(createParams)
			r := httptest.NewRequest("POST", "/tasks", reqBody)
			w := httptest.NewRecorder()

			// h.CreateTask(w, r)
			mux.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedTitle != "" {
				var task domain.Task
				_ = json.Unmarshal([]byte(w.Body.Bytes()), &task)
				assert.Equal(t, tc.expectedTitle, task.Title)
			}
		})
	}
}
