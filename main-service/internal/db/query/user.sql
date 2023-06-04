-- name: CreateUser :execresult
INSERT INTO users (
  email, nickname, password
) VALUES (
  ?, ?, ?
);

-- name: IsUserExisted :one
SELECT EXISTS(
  SELECT 1 FROM users
  WHERE email = ? LIMIT 1
);

-- name: GetUserById :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users 
SET password = ?, nickname = ?, chatbot_token = ?, updated_at = ?
WHERE id = ?;

-- name: UpdateVerificationCode :exec
UPDATE users
SET verification_code = ?, expired_at = ?, updated_at = ?
WHERE id = ?;