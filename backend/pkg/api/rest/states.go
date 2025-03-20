package rest

import (
	"context"
	"fmt"
	"net/http"
)

// statesURL – URL для получения состояний REST API.
const statesURL = "https://eaist.mos.ru/eaist2rc/api/core/states/state/list"

// FetchStates выполняет запрос для получения состояний.
func FetchStates(ctx context.Context, client *http.Client) ([]map[string]interface{}, error) {
	body := map[string]interface{}{
		"filter": map[string]interface{}{
			"categoryCode": "contractstagesupplier",
		},
	}
	items, _, err := postItems(ctx, client, statesURL, body)
	if err != nil {
		return nil, fmt.Errorf("fetch states: %v", err)
	}
	return items, nil
}
