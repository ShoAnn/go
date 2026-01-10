package domain

import (
	"context"
	"time"
)

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaskParams struct {
	Title     string
	Completed bool
}

type UpdateTaskParams struct {
	Title     string
	Completed bool
}

type TaskRepository interface {
	GetAll(ctx context.Context) ([]*Task, error)
	GetById(ctx context.Context, id int) (*Task, error)
	Create(ctx context.Context, params *CreateTaskParams) (*Task, error)
	Edit(ctx context.Context, id int) (*Task, error)
	Delete(ctx context.Context, id int) error
}
