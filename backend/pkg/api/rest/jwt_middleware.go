package rest

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ryantrue/EaistSync/pkg/config"
)

// JWTMiddleware проверяет валидность JWT токена из заголовка Authorization.
func JWTMiddleware(cfg *config.Config, log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Error("Отсутствует заголовок Authorization")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Отсутствует токен"})
			}

			// Ожидается формат: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Error("Неверный формат токена")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверный формат токена"})
			}
			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Проверка, что метод подписи корректный
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Неверный метод подписи")
				}
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil {
				log.Error("Ошибка парсинга JWT токена", zap.Error(err))
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверный токен"})
			}
			if !token.Valid {
				log.Error("Невалидный токен")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Невалидный токен"})
			}

			// Извлекаем claims и сохраняем в контекст для дальнейшего использования
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Error("Ошибка извлечения claims")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Ошибка токена"})
			}
			c.Set("user", claims)

			return next(c)
		}
	}
}
