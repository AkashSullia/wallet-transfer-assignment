package db

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type Store interface {
	Querier
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// execTx executes a function within a database transaction;
func (store *SQLStore) execTx(ctx context.Context, fn func(context.Context) error) error {
	tx, err := store.connPool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey{}, store.WithTx(tx))
	err = fn(ctx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

// serialisable

func (store *SQLStore) QueriesFrom(ctx context.Context) *Queries {
	if q, ok := ctx.Value(txKey{}).(*Queries); ok {
		return q // inside a transaction
	}
	return store.Queries // no transaction in flight
}

func RunWithSerializableRetry(
	ctx context.Context, pool *pgxpool.Pool,
	fn func(ctx context.Context, qtx *Queries) error,
) error {
	const maxAttempts = 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// 1. Begin SERIALIZABLE transaction
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: pgx.Serializable,
		})
		if err != nil {
			return err
		}

		// 2. Bind sqlc Queries to the transaction
		qtx := New(tx)

		// 3. Execute business closure
		err = fn(ctx, qtx)
		if err == nil {
			// 4. Commit
			if commitErr := tx.Commit(ctx); commitErr == nil {
				return nil
			} else {
				err = commitErr
			}
		}

		// Rollback (safe to call after commit failure)
		_ = tx.Rollback(ctx)

		// 5. Determine if retry is needed
		if !isRetryable(err) || attempt == maxAttempts {
			return err
		}

		// Exponential backoff + jitter
		backoff := time.Duration(attempt*attempt) * 25 * time.Millisecond
		jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
		select {
		case <-time.After(backoff + jitter):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("transaction failed after %d attempts", maxAttempts)
}

func isRetryable(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == "40001" || pgErr.Code == "40P01"
	}
	return false
}

func (store *SQLStore) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return RunWithSerializableRetry(ctx, store.connPool, func(ctx context.Context, qtx *Queries) error {
		ctx = context.WithValue(ctx, txKey{}, qtx)
		return fn(ctx)
	})
}
