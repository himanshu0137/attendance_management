package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

// DB is the server's standard-library database wrapper around a pgx pool.
type DB struct {
	Pool *pgxpool.Pool
	DSN  string
}

// OpenDB opens a pgx pool for the DSN and returns the wrapped database.
func OpenDB(dsn string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &DB{Pool: pool, DSN: dsn}, nil
}

// Migrate runs the goose migrations in dir against one schema, using a
// temporary pgx pool opened on the same DSN.
func (d *DB) Migrate(ctx context.Context, schema, dir string) error {
	return Migrate(ctx, schema, dir, d.DSN)
}

// Migrate runs the goose migrations in dir against one schema, using a
// temporary pgx pool opened on the same DSN.
func Migrate(ctx context.Context, schema, dir, dsn string) error {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()

	stdlib, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		return err
	}
	defer stdlib.Close()
	return goose.RunContext(ctx, "up", stdlib, dir)
}
