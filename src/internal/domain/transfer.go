package domain

import (
	"time"
	"uuid"
)

type Transfer struct {
	ID             uuid.UUID `json:"id"`
	IdempotencyKey string    `json:"idempotency_key"`
	FromWalletID   uuid.UUID `json:"from_wallet_id"`
	ToWalletID     uuid.UUID `json:"to_wallet_id"`
	Amount         int64     `json:"amount"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TransferRequest struct {
	IdempotencyKey string    `json:"idempotency_key"`
	FromWalletID   uuid.UUID `json:"from_wallet_id"`
	ToWalletID     uuid.UUID `json:"to_wallet_id"`
	Amount         int64     `json:"amount"`
}
