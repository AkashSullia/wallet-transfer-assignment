package repo

import (
	"context"
	"errors"
	"time"
	"uuid"

	db "wallet-transfer-assignment/db/generated"
	"wallet-transfer-assignment/internal/domain"

	"github.com/jackc/pgx/v5"
)

type transferRepository struct {
	store *db.SQLStore
}

func NewTransferRepository(q *db.SQLStore) domain.TransferRepository {
	return &transferRepository{q}
}

func (r *transferRepository) FindTransferByIdempotencyKey(ctx context.Context, key string) (*domain.Transfer, error) {
	q := r.store.QueriesFrom(ctx)
	tr, err := q.GetTransferByIdempotencyKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTransferNotFound
		}
		return nil, err
	}
	transfer := &domain.Transfer{
		ID:             tr.ID,
		IdempotencyKey: tr.IdempotencyKey,
		FromWalletID:   tr.FromWalletID,
		ToWalletID:     tr.FromWalletID,
		Amount:         tr.Amount,
		Status:         tr.Status,
		CreatedAt:      tr.CreatedAt,
		UpdatedAt:      tr.UpdatedAt,
	}
	return transfer, nil
}

func (r *transferRepository) Transfer(ctx context.Context, idempotencyKey string, fromWalletID, toWalletID uuid.UUID, amount int64) (*domain.Transfer, error) {
	q := r.store.QueriesFrom(ctx)
	now := time.Now()
	uuid := uuid.New()
	tr, err := q.CreateTransfer(ctx, db.CreateTransferParams{
		ID:             uuid,
		IdempotencyKey: idempotencyKey,
		FromWalletID:   fromWalletID,
		ToWalletID:     fromWalletID,
		Amount:         amount,
		Status:         "PENDING",
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		return nil, err
	}
	transfer := &domain.Transfer{
		ID:             tr.ID,
		IdempotencyKey: tr.IdempotencyKey,
		FromWalletID:   tr.FromWalletID,
		ToWalletID:     tr.FromWalletID,
		Amount:         tr.Amount,
		Status:         tr.Status,
		CreatedAt:      tr.CreatedAt,
		UpdatedAt:      tr.UpdatedAt,
	}
	return transfer, nil
}

func (r *transferRepository) UpdateTransferStatus(ctx context.Context, transferID uuid.UUID, status string) error {
	q := r.store.QueriesFrom(ctx)
	now := time.Now()
	err := q.UpdateTransferStatus(ctx, db.UpdateTransferStatusParams{
		Status:    status,
		UpdatedAt: now,
		ID:        transferID,
	})
	return err
}
