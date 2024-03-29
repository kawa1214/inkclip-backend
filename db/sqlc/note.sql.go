// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: note.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createNote = `-- name: CreateNote :one
INSERT INTO notes (
  user_id,
  title,
  content,
  is_public
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, user_id, title, content, is_public, created_at
`

type CreateNoteParams struct {
	UserID   uuid.UUID `json:"user_id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	IsPublic bool      `json:"is_public"`
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, createNote,
		arg.UserID,
		arg.Title,
		arg.Content,
		arg.IsPublic,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.IsPublic,
		&i.CreatedAt,
	)
	return i, err
}

const deleteNote = `-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1
`

func (q *Queries) DeleteNote(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteNote, id)
	return err
}

const getNote = `-- name: GetNote :one
SELECT id, user_id, title, content, is_public, created_at FROM notes
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetNote(ctx context.Context, id uuid.UUID) (Note, error) {
	row := q.db.QueryRowContext(ctx, getNote, id)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.IsPublic,
		&i.CreatedAt,
	)
	return i, err
}

const listNotesByUserId = `-- name: ListNotesByUserId :many
SELECT id, user_id, title, content, is_public, created_at FROM notes
WHERE user_id = $1
LIMIT $2
OFFSET $3
`

type ListNotesByUserIdParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) ListNotesByUserId(ctx context.Context, arg ListNotesByUserIdParams) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, listNotesByUserId, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Note{}
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.IsPublic,
			&i.CreatedAt,
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

const updateNote = `-- name: UpdateNote :one
UPDATE notes
SET
  title = $2,
  content = $3,
  is_public = $4
WHERE id = $1
RETURNING id, user_id, title, content, is_public, created_at
`

type UpdateNoteParams struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	IsPublic bool      `json:"is_public"`
}

func (q *Queries) UpdateNote(ctx context.Context, arg UpdateNoteParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, updateNote,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.IsPublic,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.IsPublic,
		&i.CreatedAt,
	)
	return i, err
}
