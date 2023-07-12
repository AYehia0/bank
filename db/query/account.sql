-- name: CreateAccount :one
INSERT INTO accounts (
  owner_name, balance, currency 
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccountById :one
SELECT * FROM accounts 
WHERE id = $1 LIMIT 1;

-- name: GetAccounts :many
SELECT * FROM accounts 
where owner_name = $1
ORDER BY id
LIMIT BY $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
