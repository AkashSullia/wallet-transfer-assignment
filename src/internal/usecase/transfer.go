package usecase

import (
	"context"
	"errors"
	"slices"
	"uuid"

	"wallet-transfer-assignment/internal/domain"
)

func NewTransferService(
	tr domain.TransferRepository,
	w domain.WalletRepository,
	le domain.LedgerEntryRepository,
	transactor domain.Transactor,
) *TransferService {
	return &TransferService{tr, w, le, transactor}
}

// Usecase layer
func (s *TransferService) Execute(ctx context.Context, req domain.TransferRequest) (transfer *domain.Transfer, err error) {
	s.transactor.WithinTx(ctx, func(ctx context.Context) error {
		transfer, err = s.doTransfer(ctx, req)
		return err
	})
	return
}

func (s *TransferService) doTransfer(ctx context.Context, req domain.TransferRequest) (*domain.Transfer, error) {
	// 1. Idempotency check
	existing, err := s.tr.FindTransferByIdempotencyKey(ctx, req.IdempotencyKey)
	if err == nil {
		if req.IdempotencyKey == existing.IdempotencyKey {
			return existing, nil
		}
	}
	if !errors.Is(err, domain.ErrTransferNotFound) {
		return nil, err
	}

	// 2. Lock wallets (sort IDs to prevent deadlocks)
	vals := []uuid.UUID{req.FromWalletID, req.ToWalletID}
	slices.SortFunc(vals, uuid.UUID.Compare)

	wallets, err := s.w.LockWallet(ctx, vals)
	if err != nil {
		return nil, err
	}

	// 3. Validate balance
	if err := domain.ValidateSufficientFunds(wallets, req); err != nil {
		return nil, err
	}

	// 4. Write transfers + ledger entries
	transfer, err := s.tr.Transfer(ctx, req.IdempotencyKey, req.FromWalletID, req.ToWalletID, req.Amount)
	if err != nil {
		return nil, err
	}
	return s.le.CreateLedgerEntries(ctx, req)
}
