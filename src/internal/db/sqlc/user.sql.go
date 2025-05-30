// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :execresult
INSERT INTO Users (
  id, name, email, password
) VALUES (
  ?, ?, ?, ?
)
`

type CreateUserParams struct {
	ID       uuid.UUID      `json:"id"`
	Name     sql.NullString `json:"name"`
	Email    sql.NullString `json:"email"`
	Password sql.NullString `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Password,
	)
}

const deleteUser = `-- name: DeleteUser :execresult
DELETE FROM Users
WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteUser, id)
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, password, created_at, updatedat FROM Users
WHERE id = ? LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.Updatedat,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, email, password, created_at, updatedat FROM Users
ORDER BY name
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.Updatedat,
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

const updateUser = `-- name: UpdateUser :execresult
UPDATE Users
SET name = ?, email = ?, password = ?
WHERE id = ?
`

type UpdateUserParams struct {
	Name     sql.NullString `json:"name"`
	Email    sql.NullString `json:"email"`
	Password sql.NullString `json:"password"`
	ID       uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUser,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.ID,
	)
}
