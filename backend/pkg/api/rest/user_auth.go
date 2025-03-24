package rest

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ryantrue/EaistSync/pkg/config"
)

// User представляет модель пользователя в БД.
type User struct {
	ID             int64     `db:"id" json:"id"`
	Username       string    `db:"username" json:"username"`
	HashedPassword string    `db:"hashed_password" json:"-"`
	Role           string    `db:"role" json:"role"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// RegisterInput описывает входные данные для регистрации.
type RegisterInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginInput описывает входные данные для авторизации.
type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse описывает структуру ответа при успешной авторизации.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// Expiry durations.
const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

// RegisterHandler обрабатывает регистрацию новых пользователей.
func RegisterHandler(cfg *config.Config, db *sqlx.DB, log *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input RegisterInput
		if err := c.Bind(&input); err != nil {
			log.Error("Ошибка парсинга данных регистрации", zap.Error(err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверные входные данные"})
		}
		if input.Username == "" || input.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username и password обязательны"})
		}

		// Проверка, существует ли уже такой пользователь
		var exists bool
		err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", input.Username)
		if err != nil && err != sql.ErrNoRows {
			log.Error("Ошибка проверки существования пользователя", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка БД"})
		}
		if exists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Пользователь с таким именем уже существует"})
		}

		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("Ошибка хеширования пароля", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка обработки данных"})
		}

		// Вставка пользователя в БД
		query := `
			INSERT INTO users (username, hashed_password)
			VALUES ($1, $2)
			RETURNING id, created_at, updated_at
		`
		var user User
		err = db.QueryRowx(query, input.Username, string(hashedPassword)).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Error("Ошибка создания пользователя", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка БД"})
		}
		user.Username = input.Username
		user.Role = "user" // по умолчанию

		return c.JSON(http.StatusCreated, user)
	}
}

// LoginHandler обрабатывает авторизацию пользователей и возвращает access и refresh токены.
func LoginHandler(cfg *config.Config, db *sqlx.DB, log *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginInput
		if err := c.Bind(&input); err != nil {
			log.Error("Ошибка парсинга данных авторизации", zap.Error(err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверные входные данные"})
		}
		if input.Username == "" || input.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username и password обязательны"})
		}

		// Поиск пользователя по имени
		var user User
		err := db.Get(&user, "SELECT * FROM users WHERE username=$1", input.Username)
		if err != nil {
			log.Error("Пользователь не найден", zap.String("username", input.Username), zap.Error(err))
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверные учетные данные"})
		}

		// Сравнение паролей
		if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password)); err != nil {
			log.Error("Неверный пароль", zap.String("username", input.Username), zap.Error(err))
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверные учетные данные"})
		}

		// Генерация access токена
		accessClaims := jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(AccessTokenExpiry).Unix(),
		}
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
		accessTokenString, err := accessToken.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			log.Error("Ошибка генерации access токена", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка авторизации"})
		}

		// Генерация refresh токена
		refreshClaims := jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(RefreshTokenExpiry).Unix(),
		}
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
		refreshTokenString, err := refreshToken.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			log.Error("Ошибка генерации refresh токена", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка авторизации"})
		}

		// Сохраняем refresh токен в хранилище
		StoreRefreshToken(refreshTokenString, user.ID)

		// Не возвращаем хеш пароля клиенту
		user.HashedPassword = ""

		return c.JSON(http.StatusOK, TokenResponse{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
			User:         user,
		})
	}
}

// RefreshTokenHandler принимает refresh токен и возвращает новый access токен (и, возможно, новый refresh токен).
func RefreshTokenHandler(cfg *config.Config, log *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload struct {
			RefreshToken string `json:"refresh_token" validate:"required"`
		}
		if err := c.Bind(&payload); err != nil {
			log.Error("Ошибка парсинга данных обновления токена", zap.Error(err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверные входные данные"})
		}
		if payload.RefreshToken == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Refresh токен обязателен"})
		}

		// Проверяем, существует ли refresh токен в нашем хранилище
		userID, ok := GetUserIDByRefreshToken(payload.RefreshToken)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверный refresh токен"})
		}

		// Парсим токен для проверки срока действия
		token, err := jwt.Parse(payload.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Неверный метод подписи")
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Невалидный refresh токен"})
		}

		// Генерируем новый access токен
		claims := jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(AccessTokenExpiry).Unix(),
		}
		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		accessTokenString, err := newAccessToken.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			log.Error("Ошибка генерации нового access токена", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка обновления токена"})
		}

		// (Опционально) можно также сгенерировать новый refresh токен,
		// обновить запись в хранилище и вернуть оба токена.
		// Для простоты здесь вернем только новый access токен.
		return c.JSON(http.StatusOK, map[string]string{"access_token": accessTokenString})
	}
}

// LogoutHandler удаляет refresh токен из хранилища.
func LogoutHandler(log *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload struct {
			RefreshToken string `json:"refresh_token" validate:"required"`
		}
		if err := c.Bind(&payload); err != nil {
			log.Error("Ошибка парсинга данных logout", zap.Error(err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверные входные данные"})
		}
		if payload.RefreshToken == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Refresh токен обязателен"})
		}

		RemoveRefreshToken(payload.RefreshToken)
		return c.JSON(http.StatusOK, map[string]string{"message": "Вы успешно вышли"})
	}
}
