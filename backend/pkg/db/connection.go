package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// NewPostgresDB создает подключение к базе данных PostgreSQL по DSN.
func NewPostgresDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// При необходимости можно добавить дополнительные настройки подключения.
	return db, nil
}
