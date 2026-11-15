-- name: CreateLedgerEntry :exec
INSERT INTO ledger_entries (id, wallet_id, transfer_id, transfer_type, amount, created_at)
VALUES ($1, $2, $3, $4, $5, $6);


