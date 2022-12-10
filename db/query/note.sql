-- name: CreateNote :one
INSERT INTO notes (
  user_id,
  title,
  content
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetNote :one
SELECT * FROM notes
WHERE id = $1 LIMIT 1;

-- name: ListNotesByUserId :many
SELECT * FROM notes
WHERE user_id = $1
LIMIT $2
OFFSET $3;

-- name: UpdateNote :one
UPDATE notes
SET
  title = $2,
  content = $3
WHERE id = $1
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1;