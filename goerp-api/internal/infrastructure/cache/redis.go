package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &redisCache{client: client}
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
