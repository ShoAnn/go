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
	Version   int32     `json:"version"`
}

type CreateTaskParams struct {
	Title     string
	Completed bool
}

type UpdateTaskParams struct {
	Title     *string // pointer here means this field is optional for updates meaning it becomes nil if not mentioned
	Completed *bool
	Version   int
}

type TaskService interface {
	ListTasks(ctx context.Context) ([]*Task, error)
	GetTask(ctx context.Context, id int) (*Task, error)
	CreateTask(ctx context.Context, p *CreateTaskParams) (*Task, error)
	CompleteTask(ctx context.Context, id int) error
	Edit(ctx context.Context, id int, p *UpdateTaskParams) (*Task, error)
	Delete(ctx context.Context, id int) error
}

type TaskRepository interface {
	GetAll(ctx context.Context) ([]*Task, error)
	GetById(ctx context.Context, id int) (*Task, error)
	Create(ctx context.Context, p *CreateTaskParams) (*Task, error)
	Update(ctx context.Context, id int, p *UpdateTaskParams) (*Task, error)
	Delete(ctx context.Context, id int) error
	MarkCompleted(ctx context.Context, id int) error
}
