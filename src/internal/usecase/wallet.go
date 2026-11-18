package usecase

import (
	"context"
	"uuid"

	"wallet-transfer-assignment/internal/domain"
)

func NewWalletService(
	w domain.WalletRepository,
) WalletService {
	return WalletService{w}
}

func (s *WalletService) CreateWallet(ctx context.Context, openingBalance int64) (*domain.Wallet, error) {
	if openingBalance < 0 {
		return nil, domain.ErrNegativeOpeningBalance
	}
	wallet, err := s.W.CreateWallet(ctx, openingBalance)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) GetWalletBalance(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	wallet, err := s.W.GetWalletBalance(ctx, id)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
