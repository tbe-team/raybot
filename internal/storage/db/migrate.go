package db

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migration/*.sql
var migrationFS embed.FS

// AutoMigrate auto migrates the database
func (p *SQLiteDB) AutoMigrate() error {
	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(p.DB, "migration"); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
