package db

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
)

// IsUniqueViolationError checks if the error is a unique violation error
func IsUniqueViolationError(err error, constraint string) bool {
	if err == nil {
		return false
	}

	var sqliteErr *sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
		return strings.Contains(sqliteErr.Error(), constraint)
	}

	return false
}

// IsNoRowsError checks if the error is a sql.ErrNoRows error
func IsNoRowsError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, sql.ErrNoRows)
}
