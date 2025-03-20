package dbutils

import (
	"encoding/base64"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// FetchRecords выполняет запрос к БД и возвращает срез записей с декодированным полем "data".
func FetchRecords(db *sqlx.DB, log *zap.Logger, query string) ([]map[string]interface{}, error) {
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			log.Error("Ошибка сканирования строки", zap.Error(err))
			continue
		}

		// Если присутствует поле "data", пытаемся его декодировать.
		if data, exists := row["data"]; exists {
			var dataStr string
			switch v := data.(type) {
			case []byte:
				dataStr = string(v)
			case string:
				dataStr = v
			}
			if dataStr != "" {
				row["data"] = decodeData(dataStr)
			}
		}

		records = append(records, row)
	}
	return records, nil
}

// decodeData пытается декодировать строку:
// сначала как base64-сериализованный JSON, затем как прямой JSON.
// Если оба метода не срабатывают, возвращается исходная строка.
func decodeData(encoded string) interface{} {
	if decodedBytes, err := base64.StdEncoding.DecodeString(encoded); err == nil {
		var decoded interface{}
		if err := json.Unmarshal(decodedBytes, &decoded); err == nil {
			return decoded
		}
	}
	var direct interface{}
	if err := json.Unmarshal([]byte(encoded), &direct); err == nil {
		return direct
	}
	return encoded
}
