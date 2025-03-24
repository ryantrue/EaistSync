package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ryantrue/EaistSync/pkg/config"
	"io"
	"net/http"
	"sort"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type apiResponse struct {
	Items []map[string]interface{} `json:"items"`
	Count int                      `json:"count"`
}

// FetchAllContracts загружает все контракты параллельно и устраняет дубликаты по полю id.
func FetchAllContracts(ctx context.Context, client *http.Client, log *zap.Logger, cfg *config.Config) ([]map[string]interface{}, error) {
	// Первый запрос для получения первой страницы и общего количества контрактов
	firstPage, totalCount, err := fetchContractsPage(ctx, client, 0, cfg.PageSize, true, cfg.ContractsURL)
	if err != nil {
		return nil, err
	}

	// Используем sync.Map для конкурентного сохранения контрактов
	var contracts sync.Map

	// Сохраняем результаты первой страницы
	for _, item := range firstPage {
		if id, ok := extractID(item); ok {
			contracts.Store(id, item)
		}
	}

	// Если страниц всего одна, возвращаем результат
	pages := (totalCount + cfg.PageSize - 1) / cfg.PageSize
	if pages <= 1 {
		return syncMapToSlice(&contracts), nil
	}

	log.Info("Всего контрактов согласно API", zap.Int("totalCount", totalCount), zap.Int("pages", pages))

	eg, ctx := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(int64(cfg.MaxConcurrency))

	// Запускаем параллельную загрузку остальных страниц
	for i := 1; i < pages; i++ {
		pageIndex := i // создаём копию переменной для замыкания
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}
		eg.Go(func() error {
			defer sem.Release(1)
			skip := pageIndex * cfg.PageSize
			pageItems, _, err := fetchContractsPage(ctx, client, skip, cfg.PageSize, false, cfg.ContractsURL)
			if err != nil {
				return fmt.Errorf("страница %d: %v", pageIndex, err)
			}
			for _, item := range pageItems {
				if id, ok := extractID(item); ok {
					contracts.Store(id, item)
				}
			}
			return nil
		})
	}

	// Ожидаем завершения всех горутин
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return syncMapToSlice(&contracts), nil
}

// fetchContractsPage выполняет запрос для получения страницы контрактов.
func fetchContractsPage(ctx context.Context, client *http.Client, skip, take int, withCount bool, contractsURL string) ([]map[string]interface{}, int, error) {
	body := buildRequestBody(skip, take, withCount)
	return postItems(ctx, client, contractsURL, body)
}

// buildRequestBody формирует тело запроса.
func buildRequestBody(skip, take int, withCount bool) map[string]interface{} {
	filter := map[string]interface{}{
		"customerId":   7884,
		"is44F3":       true,
		"is94F3":       false,
		"is223":        nil,
		"isActual":     false,
		"isOkpdChilds": false,
		"states":       []int{7, 1, 9, 5, 15, 4, 10, 3, 2, 1001, 1002, 12, 11, 5010},
	}
	return map[string]interface{}{
		"filter":    filter,
		"order":     []map[string]interface{}{{"field": "id", "desc": true}},
		"skip":      skip,
		"take":      take,
		"withCount": withCount,
	}
}

// postItems — универсальная функция для POST-запросов с JSON телом.
func postItems(ctx context.Context, client *http.Client, url string, reqBody map[string]interface{}) ([]map[string]interface{}, int, error) {
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("marshal request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, 0, fmt.Errorf("newRequest %s: %v", url, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("post %s: %v", url, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("reading response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("код=%d, тело=%s", resp.StatusCode, string(respBytes))
	}

	var result apiResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, 0, fmt.Errorf("unmarshal: %v", err)
	}
	return result.Items, result.Count, nil
}

// extractID извлекает числовой идентификатор из элемента.
func extractID(item map[string]interface{}) (int64, bool) {
	if raw, ok := item["id"]; ok {
		if f, ok := raw.(float64); ok {
			return int64(f), true
		}
	}
	return 0, false
}

// syncMapToSlice преобразует sync.Map в срез элементов и сортирует их по id.
func syncMapToSlice(m *sync.Map) []map[string]interface{} {
	var result []map[string]interface{}
	m.Range(func(_, value interface{}) bool {
		if item, ok := value.(map[string]interface{}); ok {
			result = append(result, item)
		}
		return true
	})
	sort.Slice(result, func(i, j int) bool {
		id1, _ := extractID(result[i])
		id2, _ := extractID(result[j])
		return id1 < id2
	})
	return result
}
