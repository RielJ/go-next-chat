// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  email,
  hashed_password,
  first_name,
  last_name
) VALUES (
  $1, $2, $3, $4
) RETURNING id, email, first_name, last_name, hashed_password, created_at, password_changed_at, is_email_verified
`

type CreateUserParams struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.HashedPassword,
		arg.FirstName,
		arg.LastName,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.PasswordChangedAt,
		&i.IsEmailVerified,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, email, first_name, last_name, hashed_password, created_at, password_changed_at, is_email_verified FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.PasswordChangedAt,
		&i.IsEmailVerified,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = COALESCE($1, hashed_password),
  password_changed_at = COALESCE($2, password_changed_at),
  first_name = COALESCE($3, first_name),
  last_name = COALESCE($4, last_name),
  is_email_verified = COALESCE($5, is_email_verified)
WHERE
  email = $6
RETURNING id, email, first_name, last_name, hashed_password, created_at, password_changed_at, is_email_verified
`

type UpdateUserParams struct {
	HashedPassword    pgtype.Text        `json:"hashed_password"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	FirstName         pgtype.Text        `json:"first_name"`
	LastName          pgtype.Text        `json:"last_name"`
	IsEmailVerified   pgtype.Bool        `json:"is_email_verified"`
	Email             string             `json:"email"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FirstName,
		arg.LastName,
		arg.IsEmailVerified,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.PasswordChangedAt,
		&i.IsEmailVerified,
	)
	return i, err
}
