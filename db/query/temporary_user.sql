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

-- name: GetTemporaryUserByEmailAndToken :one
SELECT * FROM temporary_users
WHERE email = $1 AND token = $2 LIMIT 1;