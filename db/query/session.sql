-- name: CreateSession :one
INSERT INTO sessions (
    id,
    username,
    refresh_token,
    is_blocked,
    ip_addr,
    user_agent,
    expired_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetSessionById :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;
