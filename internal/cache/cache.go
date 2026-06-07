package cache

import (
	"context"
	"time"
)

type Cache interface {
	GetDomainNameFromCache(ctx context.Context, key string) (string, error)
	SetDomainName(ctx context.Context, key, val string, ttl time.Duration) error
}
