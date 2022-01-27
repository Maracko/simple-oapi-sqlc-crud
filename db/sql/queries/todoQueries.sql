-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1 LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todos
ORDER BY id DESC;

-- name: ListTodosWithTags :many
SELECT * FROM todos
WHERE tags && ($1::varchar[])
ORDER BY id DESC;


-- name: CreateTodo :one
INSERT INTO todos (
  title, tags, content
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;