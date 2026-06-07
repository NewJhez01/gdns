package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func CreateNewRedisClient() *RedisClient {
	return &RedisClient{}
}

func (r *RedisClient) GetDomainNameFromCache(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to fetch domain from cache prev: %s", err)
	}
	return val, nil
}

func (r *RedisClient) SetDomainName(ctx context.Context, key, val string, ttl time.Duration) error {
	err := r.client.Set(ctx, key, val, ttl)
	if err != nil {
		return fmt.Errorf("failed to write domain name: %s into cache prev:%s", val, err)
	}
	return nil
}
