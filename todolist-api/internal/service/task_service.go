package service

import (
	db "github.com/ShoAnn/go-playground/todolist-api/internal/db/sqlc"
)

type Server struct {
	Queries *db.Queries
}
