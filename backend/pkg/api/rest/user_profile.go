package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ProfileHandler возвращает информацию о текущем пользователе,
// извлекая данные из контекста, куда middleware сохранил claims.
func ProfileHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		if user == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Неавторизованный доступ"})
		}
		return c.JSON(http.StatusOK, user)
	}
}
