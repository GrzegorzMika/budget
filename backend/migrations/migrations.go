package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5" // load postgres driver
)

const (
	MigrationsTimeout   = 5 * 60 * time.Second
	MigrationsTableName = "budget_schema_migrations"
)

// content holds our migrations content.
//
//go:embed *.sql
var content embed.FS

func PerformMigrations(ctx context.Context) error {
	_, cancel := context.WithTimeout(ctx, MigrationsTimeout)
	defer cancel()

	d, err := iofs.New(content, ".")
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: MigrationsTableName,
	})
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to perform migrations: %w", err)
	}
	return nil
}
