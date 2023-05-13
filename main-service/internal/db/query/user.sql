-- name: CreateUser :execresult
INSERT INTO users (
  username, password
) VALUES (
  ?, ?
);

-- name: IsExistUser :one
SELECT EXISTS(
  SELECT 1 FROM users
  WHERE username = ? LIMIT 1
);

-- name: GetUserById :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT ? OFFSET ?;

-- name: UpdateUser :exec
UPDATE users SET password = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;