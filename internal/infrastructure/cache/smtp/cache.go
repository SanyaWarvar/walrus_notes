package smtp

import (
	"context"
	"encoding/json"
	"wn/internal/domain/dto/auth"
	"wn/pkg/applogger"
	"wn/pkg/database/dragonfly"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	logger applogger.Logger
	client *dragonfly.Client
}

func NewCache(logger applogger.Logger, client *dragonfly.Client) *Cache {
	return &Cache{
		logger: logger,
		client: client,
	}
}

func (ch *Cache) SaveConfirmCode(ctx context.Context, email string, item auth.ConfirmationCode, ttl *time.Duration) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return ch.client.SaveValue(ctx, email, data, *ttl)
}

func (ch *Cache) GetConfirmCode(ctx context.Context, email string) (*auth.ConfirmationCode, bool, error) {
	data, err := ch.client.GetValue(ctx, email)
	if err != nil {
		switch err {
		case redis.Nil:
			return nil, false, nil

		}
		return nil, false, err
	}

	var output auth.ConfirmationCode
	err = json.Unmarshal(data, &output)
	if err != nil {
		return &output, false, err
	}
	return &output, true, nil
}
