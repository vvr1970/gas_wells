package database

import (
	"context"
	"fmt"
	"gas_wells/internal/pkg/logger"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Migrator struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

func NewMigrator(db *pgxpool.Pool, logger logger.Logger) *Migrator {
	return &Migrator{db: db, logger: logger}
}

func (m *Migrator) Run(ctx context.Context) error {
	// Создаем временную директорию для миграций
	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Записываем базовые миграции
	if err := m.createInitialMigration(tmpDir); err != nil {
		return err
	}

	// Применяем миграции
	connString := m.db.Config().ConnString()
	mig, err := migrate.New(
		"file://"+tmpDir,
		connString)
	if err != nil {
		return fmt.Errorf("failed to init migrator: %w", err)
	}

	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	m.logger.Info("Database migrations applied successfully")
	return nil
}

func (m *Migrator) createInitialMigration(dir string) error {
	// Создаем файл начальной миграции
	filePath := filepath.Join(dir, "000001_init.up.sql")
	content := `-- +migrate Up
CREATE TABLE IF NOT EXISTS wells (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	pressure FLOAT NOT NULL,
	temperature FLOAT NOT NULL,
	result FLOAT NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS wells_name_idx ON wells(name);
`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	return nil
}
