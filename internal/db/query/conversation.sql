-- name: CreateConversation :one
INSERT INTO conversation (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetConversation :one
SELECT * FROM conversation
WHERE id = $1 LIMIT 1;

-- name: UpdateConversation :one
UPDATE conversation SET
  name = $2
WHERE id = $1
RETURNING *;

-- name: GetUserConversations :many
SELECT * FROM conversation
WHERE id IN (
  SELECT conversation_id FROM conversation_users
  WHERE user_id = $1
)
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

