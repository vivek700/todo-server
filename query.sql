
-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY id;

-- name: CreateTask :one
INSERT INTO tasks (
  description, status
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateTask :exec
UPDATE tasks
set status = ?
WHERE id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?;