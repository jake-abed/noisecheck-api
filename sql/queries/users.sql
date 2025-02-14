-- name: CreateUser :one
INSERT INTO users (username, email)
  VALUES (?, ?)
  RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users ORDER BY username;
