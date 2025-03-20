package logger

import "go.uber.org/zap"

// NewLogger создаёт и возвращает новый логгер zap.
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
