package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ryantrue/EaistSync/pkg/config"
	"io"
	"net/http"
)

// Login выполняет POST-запрос для аутентификации, используя параметры из конфигурации.
func Login(ctx context.Context, client *http.Client, cfg *config.Config) error {
	body := map[string]interface{}{
		"username": cfg.Username,
		"password": cfg.Password,
		"remember": true,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal login body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", cfg.LoginURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("создание запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("чтение ответа: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("код=%d, тело=%s", resp.StatusCode, string(respBytes))
	}
	return nil
}
