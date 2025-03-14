//nolint:revive
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var _ Provider = (*DefaultProvider)(nil)

// Provider is a wrapper around *sql.DB that provides a SQLStore and a way to run transactions
type Provider interface {
	// DB returns a SQLDB
	DB() SQLDB

	// WithTX runs a function in a new transaction and rolls back if the function returns an error
	// otherwise it commits the transaction
	WithTX(ctx context.Context, fn func(SQLDB) error) error
}

type DefaultProvider struct {
	db *sql.DB
}

type Config struct {
	DBPath string
}

// NewProvider creates a new db provider
func NewProvider(cfg Config) (*DefaultProvider, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite db: %w", err)
	}

	return &DefaultProvider{db: db}, nil
}

func (p *DefaultProvider) DB() SQLDB {
	return p.db
}

func (p *DefaultProvider) WithTX(ctx context.Context, fn func(SQLDB) error) (err error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			rbErr := tx.Rollback()
			if !errors.Is(rbErr, sql.ErrTxDone) {
				err = fmt.Errorf("rollback: %w", rbErr)
			}
		}
	}()

	if err = fn(tx); err != nil {
		return fmt.Errorf("run tx: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return err
}
