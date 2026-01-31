-- name: GetAllTasks :many
SELECT * from tasks ORDER BY id;

-- name: GetTask :one
SELECT * from tasks WHERE id = $1 LIMIT 1;

-- name: CreateTask :one
INSERT INTO tasks (title, completed)
VALUES ($1, $2)
RETURNING *;

-- name: MarkCompleted :exec
UPDATE tasks
SET completed = TRUE
WHERE id = $1;

-- name: UpdateTask :one
UPDATE tasks
SET
	title = $2,
	completed = $3,
	updated_at = now(),
	version = version + 1
WHERE id = $1 AND version = $4
RETURNING *;

-- name: DeleteTask :execresult
DELETE FROM tasks
WHERE id = $1;
