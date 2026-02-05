package service

import (
	"context"
	"testing"
	"time"

	"github.com/ShoAnn/go-playground/todolist-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListTasks(t *testing.T) {
	validTasks := []*domain.Task{
		&domain.Task{
			ID:        1,
			Title:     "Shopping",
			Completed: false,
			CreatedAt: time.Now().Add(-2 * time.Minute),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		&domain.Task{
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
		expectedTasks []*domain.Task
		expectedError error
	}{
		{
			name: "Happy path",
			mockSetup: func(m *MockRepo) {
				m.On("ListTasks", mock.Anything).Return(validTasks, nil)
			},
			expectedError: nil,
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
				assert.Equal(t, scen.expectedTasks, tasks)
			}
		})
	}
}
