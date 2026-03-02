package dragonfly

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis *redis.Client
}

func New(url, username, password string) (*Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     url,
		Username: username,
		Password: password,
	})

	ping := r.Ping(context.Background())
	if ping.Err() != nil {
		return nil, ping.Err()
	}

	return &Client{redis: r}, nil
}

func (c *Client) Save(ctx context.Context, mapName, key string, payload []byte, ttl *time.Duration) error {
	err := c.redis.HSet(ctx, mapName, key, payload).Err()
	if err != nil {
		return err
	}
	if ttl != nil {
		return c.redis.HExpire(ctx, mapName, *ttl, key).Err()
	}
	return nil
}

func (c *Client) GetOne(ctx context.Context, mapName, key string) ([]byte, error) {
	return c.redis.HGet(ctx, mapName, key).Bytes()
}

func (c *Client) GetAll(ctx context.Context, mapName string) (map[string]string, error) {
	return c.redis.HGetAll(ctx, mapName).Result()
}

func (c *Client) CheckExists(ctx context.Context, mapName, key string) (bool, error) {
	return c.redis.HExists(ctx, mapName, key).Result()
}

func (c *Client) Delete(ctx context.Context, mapName string, key []string) (int64, error) {
	return c.redis.HDel(ctx, mapName, key...).Result()
}

// Close -.
func (c *Client) Close() {
	if c.redis != nil {
		c.redis.Close()
	}
}
