-- name: GetAllTasks :many
SELECT * from tasks;

-- name: GetTask :one
SELECT * from tasks WHERE id = $1;

-- name: CreateTask :one
INSERT INTO tasks (id, title, completed)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET title = $2, completed = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
