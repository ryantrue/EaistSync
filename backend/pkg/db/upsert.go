package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL драйвер
	"go.uber.org/zap"

	"github.com/ryantrue/EaistSync/pkg/utils"
)

// JSONUpserter инкапсулирует логику UPSERT‑операций для JSON-данных.
type JSONUpserter struct {
	db            *sqlx.DB
	logger        *zap.Logger
	allowedTables map[string]struct{}
}

// NewJSONUpserter создаёт новый объект JSONUpserter с динамически задаваемым списком разрешённых таблиц.
func NewJSONUpserter(db *sqlx.DB, logger *zap.Logger, allowed []string) *JSONUpserter {
	tables := make(map[string]struct{}, len(allowed))
	for _, t := range allowed {
		tables[t] = struct{}{}
	}
	return &JSONUpserter{db: db, logger: logger, allowedTables: tables}
}

// isSafeIdentifier проверяет, что имя является допустимым идентификатором.
func isSafeIdentifier(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return false
		}
	}
	return true
}

// UpsertMany выполняет транзакционное сохранение нескольких записей в указанную таблицу.
// Каждая запись сериализуется в JSON, а уникальность определяется по полю "id".
func (u *JSONUpserter) UpsertMany(ctx context.Context, table string, records []map[string]interface{}) (err error) {
	// Проверка допустимости таблицы.
	if _, ok := u.allowedTables[table]; !ok {
		err = fmt.Errorf("table %q is not allowed", table)
		u.logger.Error("UpsertMany: table not allowed", zap.String("table", table), zap.Error(err))
		return err
	}
	if !isSafeIdentifier(table) {
		err = fmt.Errorf("invalid table name %q", table)
		u.logger.Error("UpsertMany: invalid table name", zap.String("table", table), zap.Error(err))
		return err
	}

	// Если записей нет, выходим.
	if len(records) == 0 {
		u.logger.Info("No records to upsert", zap.String("table", table))
		return nil
	}

	// Контекст с таймаутом 5 секунд для транзакции.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		u.logger.Error("Failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				u.logger.Error("Rollback failed", zap.Error(rbErr))
				err = errors.Join(err, rbErr)
			} else {
				u.logger.Error("Transaction rolled back", zap.Error(err))
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				u.logger.Error("Commit failed", zap.Error(commitErr))
				if rbErr := tx.Rollback(); rbErr != nil {
					u.logger.Error("Rollback failed after commit error", zap.Error(rbErr))
					commitErr = errors.Join(commitErr, rbErr)
				}
				err = commitErr
			} else {
				u.logger.Info("Upsert successful", zap.String("table", table))
			}
		}
	}()

	// Готовим запрос UPSERT.
	query := fmt.Sprintf(`
		INSERT INTO %s (id, data)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data;
	`, table)
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement for table %s: %w", table, err)
	}
	defer stmt.Close()

	var errs []error
	for _, rec := range records {
		// Сериализация записи в JSON.
		dataBytes, jErr := json.Marshal(rec)
		if jErr != nil {
			u.logger.Warn("Error marshaling record", zap.Error(jErr), zap.Any("record", rec))
			errs = append(errs, fmt.Errorf("json marshal: %w", jErr))
			continue
		}

		// Извлечение идентификатора.
		id, idErr := utils.ExtractID(rec)
		if idErr != nil {
			u.logger.Warn("Error extracting ID", zap.Error(idErr), zap.Any("record", rec))
			errs = append(errs, idErr)
			continue
		}

		// Для каждой записи создаем отдельный контекст с таймаутом.
		execCtx, execCancel := context.WithTimeout(ctx, 5*time.Second)
		_, execErr := stmt.ExecContext(execCtx, id, dataBytes)
		execCancel()
		if execErr != nil {
			u.logger.Warn("Upsert operation failed", zap.String("table", table), zap.Error(execErr), zap.Any("id", id))
			errs = append(errs, fmt.Errorf("id %v: %w", id, execErr))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred during upsert: %w", errors.Join(errs...))
	}
	return nil
}
