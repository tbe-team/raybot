package db

import (
	"context"
	"database/sql"
)

var _ SQLDB = (*sql.DB)(nil)

// SQLDB is the interface between *sql.DB and the application
type SQLDB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
