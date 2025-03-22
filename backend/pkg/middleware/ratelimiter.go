package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients   = make(map[string]*ipLimiter)
	mu        sync.Mutex
	cleanupIn = 5 * time.Minute
)

// getLimiter возвращает лимитер для IP, создавая новый при необходимости.
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, exists := clients[ip]; exists {
		limiter.lastSeen = time.Now()
		return limiter.limiter
	}

	limiter := rate.NewLimiter(1, 5) // 1 запрос в секунду, burst до 5
	clients[ip] = &ipLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

// Cleanup удаляет лимитеры, которые не использовались более заданного времени.
func Cleanup() {
	for {
		time.Sleep(cleanupIn)
		mu.Lock()
		for ip, limiter := range clients {
			if time.Since(limiter.lastSeen) > cleanupIn {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// RateLimiterMiddleware ограничивает количество запросов от одного IP.
func RateLimiterMiddleware() echo.MiddlewareFunc {
	// Запускаем cleanup в фоне.
	go Cleanup()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := getLimiter(ip)
			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "Слишком много запросов"})
			}
			return next(c)
		}
	}
}
