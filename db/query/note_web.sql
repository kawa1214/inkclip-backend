-- name: CreateNoteWeb :one
INSERT INTO note_webs (
  note_id,
  web_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetNoteWeb :one
SELECT * FROM note_webs
WHERE note_id = $1 AND web_id = $2 LIMIT 1;

-- name: ListNoteWebsByNoteId :many
SELECT * FROM note_webs
WHERE note_id = $1;

-- name: DeleteNoteWeb :exec
DELETE FROM note_webs
WHERE note_id = $1 AND web_id = $2;

-- name: DeleteNoteWebsByNoteId :exec
DELETE FROM note_webs
WHERE note_id = $1;