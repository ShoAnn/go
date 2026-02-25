package repository

import (
	"context"

	db "github.com/ShoAnn/go/todolist-api/internal/db/sqlc"
	"github.com/ShoAnn/go/todolist-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Define an interface that both *pgxpool.Pool and pgx.Tx satisfy
type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type postgresTaskRepository struct {
	db DBTX // Change this from *pgxpool.Pool
	q  *db.Queries
}

func NewTaskRepository(conn DBTX) domain.TaskRepository {
	return &postgresTaskRepository{
		db: conn,
		q:  db.New(conn), // sqlc's New works with anything satisfying DBTX
	}
}

func (r *postgresTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	dbTasks, err := r.q.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.Task, len(dbTasks))
	for i, dbTasks := range dbTasks {
		tasks[i] = &domain.Task{
			ID:        int64(dbTasks.ID),
			Title:     dbTasks.Title,
			Completed: dbTasks.Completed,
			CreatedAt: dbTasks.CreatedAt.Time,
			UpdatedAt: dbTasks.UpdatedAt.Time,
		}
	}

	return tasks, nil
}

func (r *postgresTaskRepository) GetById(ctx context.Context, id int) (*domain.Task, error) {
	task, err := r.q.GetTask(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &domain.Task{
		ID:        int64(task.ID),
		Title:     task.Title,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt.Time,
		UpdatedAt: task.UpdatedAt.Time,
	}, nil
}

func (r *postgresTaskRepository) Create(ctx context.Context, p *domain.CreateTaskParams) (*domain.Task, error) {
	arg := &db.CreateTaskParams{
		Title:     p.Title,
		Completed: p.Completed,
	}
	task, err := r.q.CreateTask(ctx, *arg)
	if err != nil {
		return nil, err
	}

	return &domain.Task{
		ID:        int64(task.ID),
		Title:     task.Title,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt.Time,
		UpdatedAt: task.UpdatedAt.Time,
	}, nil
}

func (r *postgresTaskRepository) Update(ctx context.Context,
	id int,
	p *domain.UpdateTaskParams,
) (*domain.Task, error) {
	arg := &db.UpdateTaskParams{
		ID:        int32(id),
		Title:     *p.Title,
		Completed: *p.Completed,
	}
	task, err := r.q.UpdateTask(ctx, *arg)
	if err != nil {
		return nil, err
	}

	return &domain.Task{
		ID:        int64(task.ID),
		Title:     task.Title,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt.Time,
		UpdatedAt: task.UpdatedAt.Time,
	}, nil
}
func (r *postgresTaskRepository) Delete(ctx context.Context, id int) error {
	_, err := r.q.DeleteTask(ctx, int32(id))
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresTaskRepository) MarkCompleted(ctx context.Context, id int) error {
	err := r.q.MarkCompleted(ctx, int32(id))
	if err != nil {
		return err
	}
	return nil
}
