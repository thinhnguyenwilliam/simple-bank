-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  username,
  email,
  secret_code,
  expired_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetVerifyEmailBySecretCode :one
SELECT *
FROM verify_emails
WHERE secret_code = $1
  AND is_used = false
  AND expired_at > now();

-- name: MarkVerifyEmailUsed :exec
UPDATE verify_emails
SET is_used = true
WHERE id = $1;

-- name: VerifyUserEmailTx :exec
UPDATE users
SET is_email_verified = true
WHERE username = $1;

-- name: GetVerifyEmailByID :one
SELECT *
FROM verify_emails
WHERE id = $1;