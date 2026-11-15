-- +goose Up
CREATE TABLE  IF NOT EXISTS transfers (
  id                UUID            PRIMARY KEY,
  idempotency_key   VARCHAR(128)    NOT NULL UNIQUE,
  from_wallet_id    UUID            NOT NULL,
  to_wallet_id      UUID            NOT NULL,
  amount            BIGINT          NOT NULL,
  status            VARCHAR(32)     NOT NULL, -- PENDING, PROCESSED, FAILED
  created_at        TIMESTAMPTZ     NOT NULL,
  updated_at        TIMESTAMPTZ     NOT NULL
);

CREATE INDEX CONCURRENTLY idempotency_key_transfers ON transfers(idempotency_key);

CREATE TABLE IF NOT EXISTS ledger_entries (
  id                UUID            PRIMARY KEY,
  wallet_id         UUID            NOT NULL,
  transfer_id       UUID            NOT NULL,
  transfer_type     VARCHAR(16)     NOT NULL,
  amount            BIGINT          NOT NULL,
  created_at        TIMESTAMPTZ     NOT NULL
);

CREATE TABLE IF NOT EXISTS wallets (
  id                UUID            PRIMARY KEY,
  balance           BIGINT          NOT NULL,
  created_at        TIMESTAMPTZ     NOT NULL,
  updated_at        TIMESTAMPTZ     NOT NULL
);


-- +goose Down
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS ledger_entries;
DROP TABLE IF EXISTS wallets;
