// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, username, email)
  VALUES (?, ?, ?)
  RETURNING id, username, email, created_at, updated_at
`

type CreateUserParams struct {
	ID       string
	Username string
	Email    string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Username, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUserById = `-- name: DeleteUserById :exec
DELETE FROM users WHERE id = ?
`

func (q *Queries) DeleteUserById(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUserById, id)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, username, email, created_at, updated_at FROM users ORDER BY username ASC
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, created_at, updated_at from users WHERE lower(username) = lower(?) LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, lower string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, lower)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserById = `-- name: UpdateUserById :one
UPDATE users
  SET email = ?,
  username = ?,
  updated_at = CURRENT_TIMESTAMP
  WHERE id = ?
  RETURNING id, username, email, created_at, updated_at
`

type UpdateUserByIdParams struct {
	Email    string
	Username string
	ID       string
}

func (q *Queries) UpdateUserById(ctx context.Context, arg UpdateUserByIdParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserById, arg.Email, arg.Username, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
