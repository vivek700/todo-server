-- name: CreateUser :one
INSERT INTO users (access_code) VALUES (?) RETURNING * ;


-- name: GetUser :one
SELECT id FROM users WHERE access_code = ?;

-- name: ListTasks :many
SELECT * FROM tasks WHERE user_id = ?
ORDER BY id;

-- name: CreateTask :one
INSERT INTO tasks (
  user_id, description, status
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateTask :exec
UPDATE tasks
set status = ?
WHERE id = ? AND user_id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ? AND user_id=?;