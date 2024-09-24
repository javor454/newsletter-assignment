package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // This import is crucial
	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/logger"
)

func MigrationsUp(lg logger.Logger, pgConf *config.PostgresConfig, pgConn *sql.DB) error {
	lg.Debug("[MIGRATIONS] Starting up...")

	driver, err := postgres.WithInstance(pgConn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	if _, err := os.Stat(pgConf.MigrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", pgConf.MigrationsDir)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", pgConf.MigrationsDir),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	lg.Info("[MIGRATIONS] Done")

	return nil
}
