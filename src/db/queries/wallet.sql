-- name: CreateWallet :one
INSERT INTO wallets (id, balance, created_at, updated_at)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetWalletsForUpdate :many
SELECT * FROM wallets WHERE id = ANY($1::uuid[]) ORDER BY id FOR UPDATE;

-- name: UpdateWalletBalance :exec
UPDATE wallets SET balance = balance + $2 WHERE id = $1;

-- name: GetWalletBalance :one
SELECT * FROM wallets WHERE id = $1;

