package repo

import (
	"context"
	"time"
	"uuid"
	"wallet-transfer-assignment/internal/domain"

	db "wallet-transfer-assignment/db/generated"
)

type LedgerEntryRepository struct {
	store *db.SQLStore
}

func NewLedgerEntryRepository(q *db.SQLStore) domain.LedgerEntryRepository {
	return &LedgerEntryRepository{q}
}

func MakeLedgerEntryParam(walletID, transferID uuid.UUID, transferType string, amount int64) db.CreateLedgerEntryParams {
	now := time.Now()
	uuid := uuid.New()
	return db.CreateLedgerEntryParams{
		ID:           uuid,
		WalletID:     walletID,
		TransferID:   transferID,
		TransferType: transferType,
		Amount:       amount,
		CreatedAt:    now,
	}
}

func (r *LedgerEntryRepository) CreateLedgerEntries(ctx context.Context, tr *domain.Transfer) (*domain.Transfer, error) {
	q := r.store.QueriesFrom(ctx)
	from := MakeLedgerEntryParam(tr.FromWalletID, tr.ID, "DEBIT", tr.Amount)
	to := MakeLedgerEntryParam(tr.ToWalletID, tr.ID, "CREDIT", tr.Amount)
	for _, param := range []db.CreateLedgerEntryParams{from, to} {
		err := q.CreateLedgerEntry(ctx, param)
		if err != nil {
			return nil, err
		}
	}
	return tr, nil
}
