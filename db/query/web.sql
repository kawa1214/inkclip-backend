-- name: CreateWeb :one
INSERT INTO webs (
  user_id,
  url,
  title,
  thumbnail_url,
  html
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetWeb :one
SELECT * FROM webs
WHERE id = $1 LIMIT 1;

-- name: ListWebsByUserId :many
SELECT * FROM webs
WHERE user_id = $1
LIMIT $2
OFFSET $3;

-- name: DeleteWeb :exec
DELETE FROM webs
WHERE id = $1;

-- name: ListWebByNoteId :many
SELECT webs.* FROM webs
INNER JOIN note_webs ON webs.id = note_webs.web_id
WHERE note_webs.note_id = $1;

-- name: ListWebByNoteIds :many
SELECT webs.*, note_webs.note_id FROM webs
INNER JOIN note_webs ON webs.id = note_webs.web_id
WHERE note_webs.note_id = ANY(@ids::uuid[]);