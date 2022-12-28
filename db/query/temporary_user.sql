-- name: CreateTemporaryUser :one
INSERT INTO temporary_users (
  email,
  hashed_password,
  token,
  expires_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetTemporaryUserByToken :one
SELECT * FROM temporary_users
WHERE token = $1 LIMIT 1;