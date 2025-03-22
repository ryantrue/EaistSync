package server

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"eaistsync/backend/pkg/api/rest"
	"eaistsync/backend/pkg/config"
	"eaistsync/backend/pkg/handlers"
	"eaistsync/backend/pkg/middleware"
)

// Server хранит ссылки на базу данных, логгер и конфигурацию.
type Server struct {
	DB     *sqlx.DB
	Log    *zap.Logger
	Config *config.Config
}

// NewServer создаёт новый экземпляр Server.
func NewServer(db *sqlx.DB, log *zap.Logger, cfg *config.Config) *Server {
	return &Server{
		DB:     db,
		Log:    log,
		Config: cfg,
	}
}

// Start запускает сервер Echo на указанном адресе.
func (s *Server) Start(addr string) error {
	e := echo.New()

	// Применяем rate limiter ко всем маршрутам.
	e.Use(middleware.RateLimiterMiddleware())

	// Группа для API-эндпоинтов.
	api := e.Group("/api")

	// Существующие публичные маршруты.
	api.GET("/contracts", handlers.HandleGetRecords(s.DB, s.Log, "SELECT * FROM contracts", "contracts"))
	api.GET("/states", handlers.HandleGetRecords(s.DB, s.Log, "SELECT * FROM states", "states"))
	api.GET("/events", handlers.SSEHandler(s.Log))

	// Маршруты для регистрации и авторизации.
	api.POST("/register", rest.RegisterHandler(s.Config, s.DB, s.Log))
	api.POST("/login", rest.LoginHandler(s.Config, s.DB, s.Log))
	api.POST("/refresh", rest.RefreshTokenHandler(s.Config, s.Log))
	api.POST("/logout", rest.LogoutHandler(s.Log))

	// Защищённая группа маршрутов, требующая валидного JWT токена.
	protected := api.Group("")
	protected.Use(rest.JWTMiddleware(s.Config, s.Log))
	// Пример защищённого маршрута, возвращающего профиль пользователя.
	protected.GET("/profile", rest.ProfileHandler())

	return e.Start(addr)
}
