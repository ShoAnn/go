package service

import (
	"context"

	"github.com/ShoAnn/go-playground/todolist-api/internal/domain"
)

type TaskService struct {
	repo domain.TaskRepository
}

func (s *TaskService) GetAll(ctx context.Context) ([]*domain.Task, error) {
	taskList, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func (s *TaskService) GetById(ctx context.Context, id int) (*domain.Task, error) {
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) Create(ctx context.Context, t *domain.Task) error {
	task, err := s.repo.Create(ctx)
}

func (s *TaskService) Complete(ctx context.Context, id int) (*domain.Task, error) {}

func (s *TaskService) Edit(ctx context.Context, id int, t *domain.Task) (*domain.Task, error) {}

func (s *TaskService) Delete(ctx context.Context, id int) error {}
