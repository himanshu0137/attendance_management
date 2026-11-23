package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantTx represents a transaction with tenant isolation.
type TenantTx struct {
	tx pgx.Tx
}

// NewTenantTx creates a new TenantTx.
func NewTenantTx(tx pgx.Tx) *TenantTx {
	return &TenantTx{tx: tx}
}

// Exec executes a query within the tenant transaction.
func (t *TenantTx) Exec(ctx context.Context, sql string, arguments ...any) (pgx.CommandTag, error) {
	return t.tx.Exec(ctx, sql, arguments...)
}

// Query executes a query within the tenant transaction.
func (t *TenantTx) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	return t.tx.Query(ctx, sql, arguments...)
}

// QueryRow executes a query that returns a single row within the tenant transaction.
func (t *TenantTx) QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row {
	return t.tx.QueryRow(ctx, sql, arguments...)
}

// WithTenant wraps a function with tenant isolation.
func WithTenant(ctx context.Context, pool *pgxpool.Pool, tenantID string, fn func(context.Context, *TenantTx) error) error {
	if tenantID == "" {
		return errors.New("tenant ID cannot be empty")
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Set the tenant context for this transaction
	ctx = context.WithValue(ctx, "tenant_id", tenantID)

	tenantTx := NewTenantTx(tx)
	return fn(ctx, tenantTx)
}