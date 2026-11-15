package repo

import (
	"context"
	"uuid"

	db "wallet-transfer-assignment/db/generated"
	"wallet-transfer-assignment/internal/domain"
)

type walletRepository struct {
	store *db.SQLStore
}

func NewWalletRepository(q *db.SQLStore) domain.WalletRepository {
	return &walletRepository{q}
}

func (r *walletRepository) GetWalletBalance(ctx context.Context, walletID uuid.UUID) (*domain.Wallet, error) {
	q := r.store.QueriesFrom(ctx)
	w, err := q.GetWalletBalance(ctx, walletID)
	if err != nil {
		return nil, err
	}
	wallet := &domain.Wallet{
		ID:        w.ID,
		Balance:   w.Balance,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
	return wallet, nil
}

func (r *walletRepository) LockWallet(ctx context.Context, walletIDs []uuid.UUID) (map[uuid.UUID]*domain.Wallet, error) {
	q := r.store.QueriesFrom(ctx)
	ws, err := q.GetWalletsForUpdate(ctx, walletIDs)
	if err != nil {
		return nil, err
	}

	wallet := map[uuid.UUID]*domain.Wallet{}
	for _, w := range ws {
		wallet[w.ID] = &domain.Wallet{
			ID:        w.ID,
			Balance:   w.Balance,
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		}
	}
	return wallet, nil
}

func (r *walletRepository) UpdateWalletBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
	q := r.store.QueriesFrom(ctx)
	err := q.UpdateWalletBalance(ctx, db.UpdateWalletBalanceParams{
		ID:      walletID,
		Balance: amount,
	})
	if err != nil {
		return err
	}

	return nil
}
