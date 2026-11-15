package usecase

import (
	"context"
	"uuid"

	"wallet-transfer-assignment/internal/domain"
)

func (s *WalletService) GetWalletBalance(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	wallet, err := s.w.GetWalletBalance(ctx, id)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
