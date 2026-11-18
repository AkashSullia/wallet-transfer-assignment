package usecase

import "wallet-transfer-assignment/internal/domain"

type TransferService struct {
	tr         domain.TransferRepository
	w          domain.WalletRepository
	le         domain.LedgerEntryRepository
	transactor domain.Transactor
}

type WalletService struct {
	W domain.WalletRepository
}
