package server

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"eaistsync/backend/pkg/handlers" // Импорт пакета с обработчиками
)

// Server хранит ссылки на базу данных и логгер.
type Server struct {
	DB  *sqlx.DB
	Log *zap.Logger
}

// NewServer создаёт новый экземпляр Server.
func NewServer(db *sqlx.DB, log *zap.Logger) *Server {
	return &Server{
		DB:  db,
		Log: log,
	}
}

// Start запускает сервер Echo на указанном адресе.
func (s *Server) Start(addr string) error {
	e := echo.New()

	// Группа для API-эндпоинтов.
	api := e.Group("/api")

	// Настраиваем маршруты, передавая зависимости в обработчики.
	api.GET("/contracts", handlers.HandleGetRecords(s.DB, s.Log, "SELECT * FROM contracts", "contracts"))
	api.GET("/states", handlers.HandleGetRecords(s.DB, s.Log, "SELECT * FROM states", "states"))
	api.GET("/events", handlers.SSEHandler(s.Log))

	return e.Start(addr)
}
