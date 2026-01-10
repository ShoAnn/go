package repository

import (
	"context"

	db "github.com/ShoAnn/go-playground/todolist-api/internal/db/sqlc"
)

type DBRepo struct {
	q *db.Queries
}

func (r *DBRepo) 
