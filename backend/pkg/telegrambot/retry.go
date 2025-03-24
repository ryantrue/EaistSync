// retry.go
package telegrambot

import (
	"context"
	"fmt"
	"log"
	"time"
)

// retryOperation выполняет указанную операцию с повторными попытками.
// Если операция завершается успешно, возвращается nil. В противном случае возвращается ошибка после исчерпания попыток.
func retryOperation(ctx context.Context, operation func() error, maxRetries int, retryDelay time.Duration) error {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := operation(); err != nil {
			lastErr = err
			log.Printf("Attempt %d failed: %v", attempt+1, err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries+1, lastErr)
}
