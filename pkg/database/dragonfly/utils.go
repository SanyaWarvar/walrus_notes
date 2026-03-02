package dragonfly

import (
	"context"
	"encoding/json"
)

// Извлекает содержимое из mapName по key. Парсит json в target
func GetAndParseJson(ctx context.Context, client *Client, mapName, key string, target interface{}) error {
	rawData, err := client.GetOne(ctx, mapName, key)
	if err != nil {
		return err
	}
	if len(rawData) == 0 {
		return nil
	}
	return json.Unmarshal(rawData, target)
}
