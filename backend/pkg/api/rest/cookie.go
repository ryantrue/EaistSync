// rest/http_client.go
package rest

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/ryantrue/EaistSync/pkg/logger"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	var err error
	Logger, err = logger.NewLogger()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// NewHTTPClient создает и возвращает *http.Client с настроенным CookieJar и заданным таймаутом.
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

func NewCookieJar() (http.CookieJar, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		Logger.Error("failed to create new cookie jar", zap.Error(err))
		return nil, err
	}
	return jar, nil
}
