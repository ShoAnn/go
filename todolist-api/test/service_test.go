package test

import (
	"context"

	"github.com/ShoAnn/go-playground/todolist-api/internal/domain"
)

type MockRepo struct {
	GetAllFunc        func(ctx context.Context) ([]*domain.Task, error)
	GetByIdFunc       func(ctx context.Context, id int) (*domain.Task, error)
	CreateFunc        func(ctx context.Context, p *domain.CreateTaskParams) (*domain.Task, error)
	UpdateFunc        func(ctx context.Context, id int, p *domain.UpdateTaskParams) (*domain.Task, error)
	DeleteFunc        func(ctx context.Context, id int) error
	MarkCompletedFunc func(ctx context.Context, id int) error
}

func (m *MockRepo) GetAll(ctx context.Context) ([]*domain.Task, error) {
	return m.GetAllFunc(ctx)
}
func (m *MockRepo) GetById(ctx context.Context, id int) (*domain.Task, error) {
	return m.GetByIdFunc(ctx, id)
}
func (m *MockRepo) Create(ctx context.Context, p *domain.CreateTaskParams) (*domain.Task, error) {
	return m.CreateFunc(ctx, p)
}
func (m *MockRepo) Update(ctx context.Context, id int, p *domain.UpdateTaskParams) (*domain.Task, error) {
	return m.Update(ctx, id, p)
}
func (m *MockRepo) Delete(ctx context.Context, id int) error {
	return m.Delete(ctx, id)
}
func (m *MockRepo) MarkCompleted(ctx context.Context, id int) error {
	return m.MarkCompleted(ctx, id)
}
