-- +goose Up
-- +goose NO TRANSACTION
CREATE INDEX CONCURRENTLY idempotency_key_transfers ON transfers(idempotency_key);

-- +goose Down
-- index is deleted if table is dropped
