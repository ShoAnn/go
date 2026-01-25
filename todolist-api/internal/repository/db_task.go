package repository

import (
	"context"

	db "github.com/ShoAnn/go-playground/todolist-api/internal/db/sqlc"
)

type PostgresTaskRepository struct {
	q *db.Queries
}

func (r *PostgresTaskRepository) MarkCompleted(ctx context.Context, id int) (db.Task, error) {
	return r.q.MarkCompleted(ctx, int32(id))
}
