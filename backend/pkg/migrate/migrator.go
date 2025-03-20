package migrate

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Migrator инкапсулирует логику миграций базы данных.
type Migrator struct {
	db             *sqlx.DB
	migrationsPath string
	logger         *zap.Logger
}

// NewMigrator создает новый экземпляр Migrator.
func NewMigrator(db *sqlx.DB, migrationsPath string, logger *zap.Logger) *Migrator {
	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
		logger:         logger,
	}
}

// RunUp применяет все миграции вверх.
func (m *Migrator) RunUp() error {
	driver, err := postgres.WithInstance(m.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("не удалось создать драйвер миграций: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://"+m.migrationsPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("не удалось создать экземпляр мигратора: %w", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка при выполнении миграций: %w", err)
	}

	m.logger.Info("Миграции успешно применены")
	return nil
}
