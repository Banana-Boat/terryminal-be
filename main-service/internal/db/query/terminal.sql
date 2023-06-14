-- name: CreateTerminal :execresult
INSERT INTO terminals (
  id, name, size, remark, owner_id, template_id, total_duration
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateTerminalInfo :exec
UPDATE terminals
SET size = ?, remark = ?, total_duration = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteTerminal :exec
DELETE FROM terminals
WHERE id = ?;

-- name: GetTerminalById :one
SELECT * FROM terminals
WHERE id = ? LIMIT 1;

-- name: GetTerminalByOwnId :many
SELECT * FROM terminals
WHERE owner_id = ?;
