-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = CASE
    WHEN $1 = true THEN $2
    ELSE hashed_password
  END,
  full_name = $3,
  email = $4
WHERE username = $5
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT
  username,
  hashed_password,
  full_name,
  email,
  password_changed_at,
  created_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT
  username,
  full_name,
  email,
  created_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUserPassword :one
UPDATE users
SET
  hashed_password = $2,
  password_changed_at = now()
WHERE username = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;
