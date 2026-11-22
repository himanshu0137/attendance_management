package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// tenantCtxKey carries the tenant_id for the duration of a tenant-scoped session.
type tenantCtxKey struct{}

// WithTenant returns a context carrying tenantID.
func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantCtxKey{}, tenantID)
}

// NewTenantTx acquires a tenant-scoped connection for tenantID and returns it
// with a release function. Every query on the returned connection inherits the
// tenant context.
func NewTenantTx(ctx context.Context, pool *pgxpool.Pool, tenantID string) (*pgxpool.Conn, func(), error) {
	ctx = WithTenant(ctx, tenantID)
	c, err := pool.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	return c, func() { c.Release() }, nil
}
