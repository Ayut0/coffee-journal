package migration

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

func newMigrate(dbURL string) (*migrate.Migrate, error) {
	src, err := iofs.New(sqlFiles, "sql")
	if err != nil {
		return nil, fmt.Errorf("migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)
	if err != nil {
		return nil, fmt.Errorf("migration instance: %w", err)
	}

	return m, nil
}

func Up(dbURL string) error {
	m, err := newMigrate(dbURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange){
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

func Down(dbURL string, steps int) error {
	m, err := newMigrate(dbURL)
	if err != nil {
		return nil
	}

	if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate down: %w", err)
	}
	return nil
}