-- name: CreateUser :one
INSERT INTO users (id, username, email)
  VALUES (?, ?, ?)
  RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users ORDER BY username ASC;

-- name: GetUserByUsername :one
SELECT * from users WHERE lower(username) = lower(?) LIMIT 1;

-- name: UpdateUserById :one
UPDATE users
  SET email = ?,
  username = ?,
  updated_at = CURRENT_TIMESTAMP
  WHERE id = ?
  RETURNING *;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = ?;
