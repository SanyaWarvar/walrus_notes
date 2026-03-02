package permissions

import (
	"context"
	"encoding/json"
	"time"
	"wn/internal/domain/dto"
	"wn/internal/infrastructure/cache/common"
	"wn/pkg/applogger"
	"wn/pkg/database/dragonfly"

	"github.com/google/uuid"
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

func (ch *Cache) SavePermissionsLink(ctx context.Context, item *dto.PermissionsToken, id uuid.UUID, ttl *time.Duration) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return ch.client.Save(ctx, common.PermissionLinks, id.String(), data, ttl)
}

func (ch *Cache) GetPermissionsLink(ctx context.Context, id uuid.UUID) (*dto.PermissionsToken, bool, error) {
	data, err := ch.client.GetOne(ctx, common.PermissionLinks, id.String())
	if err != nil {
		switch err {
		case redis.Nil:
			return nil, false, nil

		}
		return nil, false, err
	}

	var output dto.PermissionsToken
	err = json.Unmarshal(data, &output)
	if err != nil {
		return &output, false, err
	}
	return &output, true, nil
}
