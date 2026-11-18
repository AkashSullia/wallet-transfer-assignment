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
) TransferService {
	return TransferService{tr, w, le, transactor}
}

func (s *TransferService) Transfer(ctx context.Context, req domain.TransferReq) (transfer *domain.Transfer, err error) {
	txErr := s.transactor.WithinTx(ctx, func(ctx context.Context) error {
		transfer, err = s.doTransfer(ctx, req)
		return err
	})
	if txErr != nil {
		return nil, txErr
	}
	return
}

func (s *TransferService) doTransfer(ctx context.Context, req domain.TransferReq) (*domain.Transfer, error) {
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

	// 3. Initiate transfer
	transfer, err := s.tr.Transfer(ctx, req.IdempotencyKey, req.FromWalletID, req.ToWalletID, req.Amount)
	if err != nil {
		return nil, err
	}

	// 4. Validate balance, mark transfer as failed if check fails and return
	if err := domain.ValidateSufficientFunds(wallets, req); err != nil {
		if updErr := s.tr.UpdateTransferStatus(ctx, transfer.ID, "FAILED"); updErr != nil {
			return nil, err
		}
		return transfer, err
	}

	// 5. Create ledger entries
	_, err = s.le.CreateLedgerEntries(ctx, transfer)
	if err != nil {
		if updErr := s.tr.UpdateTransferStatus(ctx, transfer.ID, "FAILED"); updErr != nil {
			return nil, err
		}
		return transfer, nil
	}

	// 6. Update final transfer status to succesful
	err = s.tr.UpdateTransferStatus(ctx, transfer.ID, "PROCESSED")
	if err != nil {
		return nil, err
	}
	return transfer, nil
}
