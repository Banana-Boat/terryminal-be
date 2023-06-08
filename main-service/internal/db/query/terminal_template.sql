-- name: CreateTerminalTemplate :execresult
INSERT INTO terminal_templates (
  name, image_name, size, description
) VALUES (
  ?, ?, ?, ?
);

-- name: GetTerminalTemplates :many
SELECT * FROM terminal_templates;

-- name: GetTerminalTemplateById :one
SELECT * FROM terminal_templates
WHERE id = ? LIMIT 1;