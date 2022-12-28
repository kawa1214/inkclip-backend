// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: temporary_user.sql

package db

import (
	"context"
	"time"
)

const createTemporaryUser = `-- name: CreateTemporaryUser :one
INSERT INTO temporary_users (
  email,
  hashed_password,
  token,
  expires_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING email, hashed_password, token, expires_at, created_at
`

type CreateTemporaryUserParams struct {
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Token          string    `json:"token"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func (q *Queries) CreateTemporaryUser(ctx context.Context, arg CreateTemporaryUserParams) (TemporaryUser, error) {
	row := q.db.QueryRowContext(ctx, createTemporaryUser,
		arg.Email,
		arg.HashedPassword,
		arg.Token,
		arg.ExpiresAt,
	)
	var i TemporaryUser
	err := row.Scan(
		&i.Email,
		&i.HashedPassword,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getTemporaryUserByToken = `-- name: GetTemporaryUserByToken :one
SELECT email, hashed_password, token, expires_at, created_at FROM temporary_users
WHERE token = $1 LIMIT 1
`

func (q *Queries) GetTemporaryUserByToken(ctx context.Context, token string) (TemporaryUser, error) {
	row := q.db.QueryRowContext(ctx, getTemporaryUserByToken, token)
	var i TemporaryUser
	err := row.Scan(
		&i.Email,
		&i.HashedPassword,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}