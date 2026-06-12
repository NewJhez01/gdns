package cache

import (
	"context"
	"errors"
	"time"
)

type Value struct {
	Rdata     []byte
	IsBlocked bool
	TTL       uint32
}

var ErrEmpty = errors.New("no val in cache")

type Cache interface {
	GetDomainNameFromCache(ctx context.Context, key string) (Value, error)
	SetDomainName(ctx context.Context, key string, v Value, ttl time.Duration) error
}
