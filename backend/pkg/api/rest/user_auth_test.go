package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ryantrue/EaistSync/pkg/config"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func TestLoginHandler(t *testing.T) {
	// Создаем поддельную БД и mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Создаем тестовый логгер
	logger, _ := zap.NewDevelopment()

	// Создаем минимальную конфигурацию для теста
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	// Подготавливаем mock: при запросе пользователя с именем "test", возвращаем фиктивные данные.
	rows := sqlmock.NewRows([]string{"id", "username", "hashed_password", "role", "created_at", "updated_at"}).
		AddRow(1, "test", "$2a$10$CwTycUXWue0Thq9StjUM0uJ8OFY8uY1k08g2z9vZFK/2dYY/PkB2K", "user", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM users WHERE username=\\$1").
		WithArgs("test").
		WillReturnRows(rows)

	// Создаем echo context с тестовым запросом
	e := echo.New()
	loginPayload := LoginInput{
		Username: "test",
		Password: "password", // пароль должен совпадать с bcrypt-значением; для теста можно изменить хеш
	}
	body, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	handler := LoginHandler(cfg, sqlxDB, logger)
	if err := handler(c); err != nil {
		t.Fatalf("Обработчик вернул ошибку: %v", err)
	}

	// Проверяем код ответа
	if rec.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", rec.Code)
	}

	// Можно дополнительно проверить содержимое ответа, распарсив JSON.
	var resp TokenResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Ошибка парсинга ответа: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Error("Токены не возвращены")
	}
}
