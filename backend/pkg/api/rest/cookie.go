package rest

import (
	"eaistsync/backend/pkg/logger"
	"net/http"
	"net/http/cookiejar"
	"time"

	"go.uber.org/zap"
)

var Logger *zap.Logger

// init инициализирует глобальный логгер, используя функцию NewLogger из пакета logger.
func init() {
	var err error
	Logger, err = logger.NewLogger()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// NewHTTPClient создаёт HTTP-клиент с CookieJar и заданным таймаутом.
// При возникновении ошибки при создании CookieJar, ошибка логируется и возвращается.
func NewHTTPClient(timeout time.Duration) (*http.Client, error) {
	jar, err := NewCookieJar()
	if err != nil {
		Logger.Error("failed to create cookie jar", zap.Error(err))
		return nil, err
	}
	return &http.Client{
		Jar:     jar,
		Timeout: timeout,
	}, nil
}

// NewCookieJar создаёт и возвращает новый CookieJar.
// Ошибки при создании CookieJar логируются и возвращаются вызывающей стороне.
func NewCookieJar() (http.CookieJar, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		Logger.Error("failed to create new cookie jar", zap.Error(err))
		return nil, err
	}
	return jar, nil
}
