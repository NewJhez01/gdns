package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func CreateNewRedisClient(c *redis.Client) *RedisClient {
	return &RedisClient{
		client: c,
	}
}

func (r *RedisClient) GetDomainNameFromCache(ctx context.Context, key string) (Value, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return Value{}, ErrEmpty
	}
	if err != nil {
		return Value{}, fmt.Errorf("failed to fetch domain from cache prev: %s", err)
	}
	val := Value{}
	err = json.Unmarshal([]byte(res), &val)
	if err != nil {
		_ = r.client.Del(ctx, key)
		return Value{}, ErrEmpty
	}
	return val, nil
}

func (r *RedisClient) SetDomainName(ctx context.Context, key string, val Value, ttl time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, b, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to write domain name: %s into cache prev:%s", val.Answer, err)
	}
	return nil
}
