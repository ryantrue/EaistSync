package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"eaistsync/backend/pkg/dbutils" // Импорт пакета с утилитами для работы с БД
)

// HandleGetRecords возвращает обработчик для GET-запросов, который выбирает данные по указанному запросу.
func HandleGetRecords(db *sqlx.DB, log *zap.Logger, query string, sourceName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		records, err := dbutils.FetchRecords(db, log, query)
		if err != nil {
			log.Error("Ошибка получения данных для "+sourceName, zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения данных"})
		}
		return c.JSON(http.StatusOK, records)
	}
}

// SSEHandler реализует сервер-сент эвенты (SSE) по адресу /api/events.
func SSEHandler(log *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := c.Response()
		res.Header().Set("Content-Type", "text/event-stream")
		res.Header().Set("Cache-Control", "no-cache")
		res.Header().Set("Connection", "keep-alive")

		sendEvent := func(data string) {
			fmt.Fprintf(res, "data: %s\n\n", data)
			res.Flush()
		}

		// Отправляем событие сразу при подключении.
		sendEvent(fmt.Sprintf("Событие: %s", time.Now().Format(time.RFC1123)))

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-c.Request().Context().Done():
				return nil
			case t := <-ticker.C:
				sendEvent(fmt.Sprintf("Событие: %s", t.Format(time.RFC1123)))
			}
		}
	}
}
