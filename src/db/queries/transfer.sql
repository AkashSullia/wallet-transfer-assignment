-- name: CreateTransfer :one
INSERT INTO transfers (id, idempotency_key, from_wallet_id, to_wallet_id, amount, status, created_at, updated_at)
VALUES ( $1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateTransferStatus :exec
UPDATE transfers
SET status = $1, updated_at = $2
WHERE id = $3;

-- name: GetTransferByIdempotencyKey :one
SELECT * FROM transfers WHERE idempotency_key = $1;
