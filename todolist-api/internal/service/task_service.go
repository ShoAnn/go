package service

import (
	"context"

	"github.com/ShoAnn/go-playground/todolist-api/internal/domain"
)

type TaskService struct {
	repo domain.TaskRepository
}

func (s *TaskService) ListTasks(ctx context.Context) ([]*domain.Task, error) {
	taskList, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func (s *TaskService) GetTask(ctx context.Context, id int) (*domain.Task, error) {
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) CreateTask(ctx context.Context, p *domain.CreateTaskParams) (*domain.Task, error) {
	task, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return &domain.Task{
		Title:     task.Title,
		Completed: task.Completed,
	}, nil
}

func (s *TaskService) CompleteTask(ctx context.Context, id int) (*domain.Task, error) {
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	params := &domain.UpdateTaskParams{
		Title:     task.Title,
		Completed: true,
	}
	newTask, err := s.repo.Update(ctx, id, params)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}

func (s *TaskService) Edit(ctx context.Context, id int, t *domain.Task) (*domain.Task, error) {}

func (s *TaskService) Delete(ctx context.Context, id int) error {}
