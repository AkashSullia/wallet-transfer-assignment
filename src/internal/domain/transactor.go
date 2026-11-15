package domain

import "context"

type Transactor interface {
	// WithinTx executes fn within a transaction. The implementation decides the isolation level (SERIALIZABLE)
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
