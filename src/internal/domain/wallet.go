package domain

import (
	"errors"
	"math"
	"time"
	"uuid"
)

type Wallet struct {
	ID        uuid.UUID `json:"id"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

var (
	ErrInsufficientFunds = errors.New("error insufficient fund")
	ErrFundOverflow      = errors.New("error fund overflow")
)

func ValidateSufficientFunds(wallets map[uuid.UUID]*Wallet, req TransferRequest) error {
	fromWallet := wallets[req.FromWalletID]
	toWallet := wallets[req.ToWalletID]
	if err := fromWallet.DebitCheck(req.Amount); err != nil {
		return err
	}
	if err := toWallet.CreditCheck(req.Amount); err != nil {
		return err
	}
	return nil
}

func (w *Wallet) CreditCheck(amount int64) error {
	if amount > 0 && w.Balance > math.MaxInt64-amount {
		return ErrFundOverflow
	}
	return nil
}

func (w *Wallet) DebitCheck(amount int64) error {
	if w.Balance-amount <= 0 {
		return ErrInsufficientFunds
	}
	return nil
}
