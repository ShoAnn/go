-- name: GetAllTasks :many
SELECT * from tasks ORDER BY id;

-- name: GetTask :one
SELECT * from tasks WHERE id = $1 LIMIT 1;

-- name: CreateTask :one
INSERT INTO tasks (title, completed)
VALUES ($1, $2)
RETURNING *;

-- name: MarkCompleted :one
UPDATE tasks
SET completed = TRUE
WHERE id = $1
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET title = $2, completed = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTask :execresult
DELETE FROM tasks
WHERE id = $1;
