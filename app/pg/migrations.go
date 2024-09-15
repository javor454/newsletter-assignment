package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // This import is crucial
)

func MigrationsUp(pgConn *sql.DB) error {
	driver, err := postgres.WithInstance(pgConn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %s", err.Error())
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %s", err.Error())
	}

	// Construct the absolute path to the migrations directory
	migrationsPath := filepath.Join(cwd, "migrations")

	// Check if the migrations directory exists
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationsPath)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %s", err.Error())
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %s", err.Error())
	}

	return nil
}
