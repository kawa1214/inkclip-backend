-- name: CreateWeb :one
INSERT INTO webs (
  user_id,
  url,
  title,
  thumbnail_url
) VALUES (
  $1, $2, $3, $4
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
