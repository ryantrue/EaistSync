package utils

import (
	"fmt"
	"strconv"
)

// ExtractID извлекает числовой идентификатор из элемента.
// Поддерживаемые типы: int, int64, float64, string (строка парсится как десятичное число).
// При отсутствии поля "id" или неподдерживаемом типе возвращается ошибка.
func ExtractID(item map[string]interface{}) (int64, error) {
	raw, ok := item["id"]
	if !ok {
		return 0, fmt.Errorf("ключ 'id' не найден")
	}

	switch v := raw.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("не удалось преобразовать строку %q в int64: %w", v, err)
		}
		return id, nil
	default:
		return 0, fmt.Errorf("неподдерживаемый тип для поля id: %T", raw)
	}
}
