package utils

import (
	"fmt"

	"github.com/spf13/cast"
)

// ExtractID извлекает числовой идентификатор из элемента.
// Поддерживаются типы: int, int64, float64, string (строка парсится как десятичное число).
// При отсутствии поля "id" или неподдерживаемом типе возвращается ошибка.
func ExtractID(item map[string]interface{}) (int64, error) {
	raw, ok := item["id"]
	if !ok {
		return 0, fmt.Errorf("ключ 'id' не найден")
	}

	// Преобразование значения в int64 с использованием пакета cast
	id, err := cast.ToInt64E(raw)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать значение %v в int64: %w", raw, err)
	}

	return id, nil
}
