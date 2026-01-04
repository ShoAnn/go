package domain

import "time"

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TodoRepository interface {
	GetAll(task *Task) ([]*Task, error)
	GetById(id int) (*Task, error)
	Create(task *Task) error
	Update(id int) (*Task, error)
	Delete(id int) error
}
