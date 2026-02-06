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
				m.On("GetAll", mock.Anything).Return(nil, errors.New("db failure"))
			},
			expectedError: errors.New("db failure"),
		},
	}
	for _, scen := range tests {
		t.Run(scen.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			scen.mockSetup(mockRepo)

			svc := &TaskService{repo: mockRepo}

			tasks, err := svc.ListTasks(context.Background())

			if scen.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, tasks, scen.expectedCount)
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
			name:   "invalid id",
			taskId: -1,
			mockSetup: func(m *MockRepo) {
				m.On("GetById", mock.Anything, -1).Return(nil, errors.New("invalid id"))
			},
			expectedError: errors.New("invalid id"),
		},
	}

	for _, scen := range tests {
		t.Run(scen.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			scen.mockSetup(mockRepo)

			svc := &TaskService{repo: mockRepo}

			task, err := svc.GetTask(context.Background(), scen.taskId)
			if scen.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, task.Title, scen.expectedTitle)
			}
		})
	}
}
