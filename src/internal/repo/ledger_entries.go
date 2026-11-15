package repo

import (
	"context"

	db "wallet-transfer-assignment/db/generated"
	"wallet-transfer-assignment/internal/domain"
)

type LedgerEntryRepository struct {
	store *db.SQLStore
}

func NewLedgerEntryRepository(q *db.SQLStore) domain.LedgerEntryRepository {
	return &LedgerEntryRepository{q}
}

func (r *LedgerEntryRepository) CreateLedgerEntries(ctx context.Context, req domain.TransferRequest) (*domain.Transfer, error) {
	return nil, nil
}
