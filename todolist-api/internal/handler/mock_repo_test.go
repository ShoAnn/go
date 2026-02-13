package handler

import (
	"context"

	"github.com/ShoAnn/go-playground/todolist-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) ListTasks(ctx context.Context) ([]*domain.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}
func (m *MockTaskService) GetTask(ctx context.Context, id int) (*domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}
func (m *MockTaskService) CreateTask(ctx context.Context, p *domain.CreateTaskParams) (*domain.Task, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}
func (m *MockTaskService) EditTask(ctx context.Context, id int, p *domain.UpdateTaskParams) (*domain.Task, error) {
	args := m.Called(ctx, id, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}
func (m *MockTaskService) CompleteTask(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Error(0)
}
func (m *MockTaskService) DeleteTask(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Error(0)
}
