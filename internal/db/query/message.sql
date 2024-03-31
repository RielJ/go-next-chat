-- name: CreateMessage :one
INSERT INTO message (
  user_id,
  message,
  conversation_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetMessage :one
SELECT * FROM message
WHERE id = $1 LIMIT 1;

-- name: UpdateMessage :one
UPDATE message
SET
  message = COALESCE(sqlc.narg(message), message)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: GetConversationMessages :many
SELECT * FROM message
WHERE conversation_id = $1
AND user_id = $2
ORDER BY created_at DESC
LIMIT $3
OFFSET $4;


