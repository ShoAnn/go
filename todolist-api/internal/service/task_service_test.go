package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ShoAnn/go-playground/todolist-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListTasks(t *testing.T) {
	validTasks := []*domain.Task{
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

	tests := []struct {
		name          string
		mockSetup     func(m *MockRepo)
		expectedCount int
		expectedError error
	}{
		{
			name: "Happy path",
			mockSetup: func(m *MockRepo) {
				m.On("GetAll", mock.Anything).Return(validTasks, nil)
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "empty data",
			mockSetup: func(m *MockRepo) {
				m.On("GetAll", mock.Anything).Return([]*domain.Task{}, nil)
			},
			expectedCount: 0,
		},
		{
			name: "db failure",
			mockSetup: func(m *MockRepo) {
				m.On("GetAll", mock.Anything).Return(nil, errors.New("repo error"))
			},
			expectedError: errors.New("repo error"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tc.mockSetup(mockRepo)

			svc := &TaskService{repo: mockRepo}

			tasks, err := svc.ListTasks(context.Background())

			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, tasks, tc.expectedCount)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	validTask := &domain.Task{
		ID:        1,
		Title:     "Shopping",
		Completed: false,
		CreatedAt: time.Now().Add(-2 * time.Minute),
		UpdatedAt: time.Now(),
		Version:   0,
	}

	tests := []struct {
		name          string
		taskId        int
		mockSetup     func(m *MockRepo)
		expectedTitle string
		expectedError error
	}{
		{
			name:   "Happy",
			taskId: 1,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, 1).Return(validTask, nil)
			},
			expectedTitle: "Shopping",
		},
		{
			name:   "id not found",
			taskId: 2,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, 2).Return(nil, errors.New("id not found"))
			},
			expectedError: errors.New("id not found"),
		},
		{
			name:   "db failure",
			taskId: 1,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, 1).Return(nil, errors.New("repo error"))
			},
			expectedError: errors.New("repo error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tc.mockSetup(mockRepo)

			svc := &TaskService{repo: mockRepo}

			task, err := svc.GetTask(context.Background(), tc.taskId)
			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, task.Title, tc.expectedTitle)
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	validTask := &domain.Task{
		ID:        1,
		Title:     "Shopping",
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   0,
	}
	params := &domain.CreateTaskParams{Title: "Shopping", Completed: false}

	tests := []struct {
		name          string
		params        *domain.CreateTaskParams
		mockSetup     func(m *MockRepo)
		expectedTitle string
		expectedError error
	}{
		{
			name:   "Happy",
			params: params,
			mockSetup: func(m *MockRepo) {
				m.On("Create", mock.Anything, params).Return(validTask, nil)
			},
			expectedTitle: "Shopping",
			expectedError: nil,
		},
		{
			name:   "repo error",
			params: params,
			mockSetup: func(m *MockRepo) {
				m.On("Create", mock.Anything, params).Return(nil, errors.New("create task failed"))
			},
			expectedTitle: "",
			expectedError: errors.New("create task failed"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tc.mockSetup(mockRepo)

			svc := &TaskService{repo: mockRepo}

			task, err := svc.CreateTask(context.Background(), tc.params)

			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, task.Title, tc.expectedTitle)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	now := time.Now()
	task := &domain.Task{
		ID:        1,
		Title:     "Shopping",
		Completed: false,
		CreatedAt: now,
		UpdatedAt: time.Now(),
		Version:   0,
	}

	newTitle := "Shopping Grocery"
	newCompleted := true
	params := &domain.UpdateTaskParams{Title: &newTitle, Completed: &newCompleted, Version: 0}

	updatedTask := &domain.Task{
		ID:        1,
		Title:     newTitle,
		Completed: newCompleted,
		CreatedAt: now,
		UpdatedAt: time.Now(),
		Version:   0,
	}

	tests := []struct {
		name                 string
		taskId               int
		params               *domain.UpdateTaskParams
		mockSetup            func(m *MockRepo)
		expectedNewTitle     string
		expectedNewCompleted bool
		expectedError        error
	}{
		{
			name:   "Happy",
			taskId: 1,
			params: params,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, 1).Return(task, nil)
				m.On("Update", mock.Anything, 1, params).Return(updatedTask, nil)
			},
			expectedNewTitle:     "Shopping Grocery",
			expectedNewCompleted: true,
			expectedError:        nil,
		},
		{
			name:   "invalid id",
			taskId: -1,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, -1).Return(nil, errors.New("invalid id"))
			},
			expectedNewTitle:     "",
			expectedNewCompleted: false,
			expectedError:        errors.New("invalid id"),
		},
		{
			name:   "task not found",
			taskId: 999,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, 999).Return(nil, errors.New("invalid id"))
			},
			expectedNewTitle:     "",
			expectedNewCompleted: false,
			expectedError:        errors.New("invalid id"),
		},
		{
			name:   "repo update failed",
			taskId: 1,
			params: params,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, 1).Return(task, nil)
				m.On("Update", mock.Anything, 1, params).Return(nil, errors.New("repo error"))
			},
			expectedNewTitle:     "",
			expectedNewCompleted: false,
			expectedError:        errors.New("repo error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tc.mockSetup(mockRepo)

			svc := &TaskService{repo: mockRepo}

			task, err := svc.EditTask(context.Background(), tc.taskId, tc.params)

			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, task.Title, tc.expectedNewTitle)
				assert.Equal(t, task.Completed, tc.expectedNewCompleted)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name          string
		taskId        int
		mockSetup     func(m *MockRepo)
		expectedError error
	}{
		{
			name:   "Happy path",
			taskId: 1,
			mockSetup: func(m *MockRepo) {
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "task not founc",
			taskId: 99,
			mockSetup: func(m *MockRepo) {
				m.On("Delete", mock.Anything, 99).Return(errors.New("id not found"))
			},
			expectedError: errors.New("id not found"),
		},
		{
			name:   "repo failure",
			taskId: 1,
			mockSetup: func(m *MockRepo) {
				m.On("Delete", mock.Anything, 1).Return(errors.New("repo error"))
			},
			expectedError: errors.New("repo error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tc.mockSetup(mockRepo)

			svc := &TaskService{repo: mockRepo}

			err := svc.DeleteTask(context.Background(), tc.taskId)
			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMarkCompleted(t *testing.T) {
	incompleteTask := &domain.Task{
		ID:        1,
		Title:     "Shopping",
		Completed: false,
		CreatedAt: time.Now().Add(-2 * time.Minute),
		UpdatedAt: time.Now(),
		Version:   0,
	}
	completedTask := &domain.Task{
		ID:        2,
		Title:     "Exercise",
		Completed: true,
		CreatedAt: time.Now().Add(-2 * time.Minute),
		UpdatedAt: time.Now(),
		Version:   0,
	}

	tests := []struct {
		name              string
		taskId            int
		mockSetup         func(m *MockRepo)
		expectedCompleted bool
		expectedError     error
	}{
		{
			name:   "Happy path",
			taskId: 1,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, 1).Return(incompleteTask, nil)
				m.On("MarkCompleted", mock.Anything, 1).Return(nil)
			},
			expectedCompleted: true,
			expectedError:     nil,
		},
		{
			name:   "task not found",
			taskId: 99,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, 99).Return(nil, nil)
				m.On("MarkCompleted", mock.Anything, 99).Return(errors.New("id not found"))
			},
			expectedCompleted: false,
			expectedError:     errors.New("id not found"),
		},
		{
			name:   "task already completed",
			taskId: 2,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, 2).Return(completedTask, nil)
				m.On("MarkCompleted", mock.Anything, 2).Return(errors.New("task already completed"))
			},
			expectedCompleted: true,
			expectedError:     errors.New("task already completed"),
		},
		{
			name:   "repo failure",
			taskId: 1,
			mockSetup: func(m *MockRepo) {
				m.On("MarkCompleted", mock.Anything, 1).Return(errors.New("repo error"))
			},
			expectedCompleted: false,
			expectedError:     errors.New("repo error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tc.mockSetup(mockRepo)

			svc := &TaskService{repo: mockRepo}

			err := svc.repo.MarkCompleted(context.Background(), tc.taskId)
			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
