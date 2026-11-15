package domain

import (
	"context"
	"errors"
	"uuid"
)

type TxManager interface {
	execTx(ctx context.Context, fn func(context.Context) error) error
}

type TransferRepository interface {
	Transfer(ctx context.Context, idempotencyKey string, fromWalletID, toWalletID uuid.UUID, amount int64) (*Transfer, error)
	FindTransferByIdempotencyKey(ctx context.Context, key string) (*Transfer, error)
}

type WalletRepository interface {
	GetWalletBalance(ctx context.Context, walletID uuid.UUID) (*Wallet, error)
	LockWallet(ctx context.Context, walletIDs []uuid.UUID) (map[uuid.UUID]*Wallet, error)
	UpdateWalletBalance(ctx context.Context, walletID uuid.UUID, amount int64) error
}

type LedgerEntryRepository interface {
	CreateLedgerEntries(ctx context.Context, req TransferRequest) (*Transfer, error)
}

var ErrTransferNotFound = errors.New("record not found")
