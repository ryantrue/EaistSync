package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"eaistsync/backend/pkg/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL драйвер
	"go.uber.org/zap"
)

// allowedTables задаёт список разрешённых таблиц для операций с БД.
var allowedTables = map[string]struct{}{
	"contracts": {},
	"states":    {},
	"lots":      {},
}

// isValidTable проверяет, разрешено ли использование указанной таблицы.
func isValidTable(tableName string) bool {
	_, ok := allowedTables[strings.ToLower(tableName)]
	return ok
}

// UpsertJSON сохраняет данные в указанную таблицу с использованием UPSERT'а.
// Данные сохраняются в формате JSONB (в столбце "data") с уникальным идентификатором "id".
// Операция выполняется внутри транзакции для обеспечения атомарности.
func UpsertJSON(ctx context.Context, db *sqlx.DB, tableName string, items []map[string]interface{}, logger *zap.Logger) error {
	if len(items) == 0 {
		logger.Info("Нет данных для сохранения", zap.String("table", tableName))
		return nil
	}

	if !isValidTable(tableName) {
		return fmt.Errorf("таблица %s не разрешена", tableName)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, data)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data;
	`, tableName)

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer func() {
		// Если ошибка возникла — откатываем транзакцию.
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Error("Ошибка отката транзакции", zap.Error(rbErr))
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("не удалось подготовить запрос для таблицы %s: %w", tableName, err)
	}
	defer stmt.Close()

	var errs []error
	for _, item := range items {
		dataBytes, err := json.Marshal(item)
		if err != nil {
			logger.Warn("Ошибка сериализации элемента", zap.Error(err), zap.Any("item", item))
			errs = append(errs, fmt.Errorf("сериализация: %w", err))
			continue
		}

		id, err := utils.ExtractID(item)
		if err != nil {
			logger.Warn("Ошибка извлечения ID", zap.Error(err), zap.Any("item", item))
			errs = append(errs, err)
			continue
		}

		// Создаем контекст с таймаутом для каждого запроса.
		execCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, execErr := stmt.ExecContext(execCtx, id, dataBytes)
		cancel()
		if execErr != nil {
			logger.Warn("Ошибка UPSERT операции", zap.String("table", tableName), zap.Error(execErr))
			errs = append(errs, fmt.Errorf("id %v: %w", id, execErr))
		}
	}

	if len(errs) > 0 {
		// Для Go 1.20+ можно использовать errors.Join(errs...)
		return fmt.Errorf("ошибки при сохранении данных: %w", errors.Join(errs...))
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}
	return nil
}
