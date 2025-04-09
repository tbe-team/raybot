package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type Provider interface {
	DB
	WithTX(ctx context.Context, fn func(DB) error) error
}

var _ Provider = (*SQLiteDB)(nil)

type SQLiteDB struct {
	*sql.DB
}

func NewSQLiteDB(path string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return &SQLiteDB{DB: db}, nil
}

func (p *SQLiteDB) WithTX(ctx context.Context, fn func(DB) error) (err error) {
	tx, err := p.DB.BeginTx(ctx, nil)
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
