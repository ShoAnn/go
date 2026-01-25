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
	Title     *string // pointer here means this field is optional for updates meaning it becomes nil if not mentioned
	Completed *bool
}

type TaskRepository interface {
	GetAll(ctx context.Context) ([]*Task, error)
	GetById(ctx context.Context, id int) (*Task, error)
	Create(ctx context.Context, params *CreateTaskParams) (*Task, error)
	MarkCompleted(ctx context.Context, id int) (*Task, error)
	Update(ctx context.Context, id int, params *UpdateTaskParams) (*Task, error)
	Delete(ctx context.Context, id int) error
}
