-- name: CreateConversationUser :one
INSERT INTO conversation_users (
  conversation_id,
  user_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetConversationUser :one
SELECT * FROM conversation_users
WHERE id = $1 LIMIT 1;

-- name: DeleteConversationUser :one
DELETE FROM conversation_users
WHERE user_id = $1
AND conversation_id = $2
RETURNING *;

-- name: GetConversationUsers :many
SELECT * FROM conversation_users
WHERE conversation_id = $1;
