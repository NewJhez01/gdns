package cache

import (
	"context"
	"errors"
	"time"
)

type Value struct {
	Answer    []byte
	IsBlocked bool
}

var ErrEmpty = errors.New("no val in cache")

type Cache interface {
	GetDomainNameFromCache(ctx context.Context, key string) (Value, error)
	SetDomainName(ctx context.Context, key string, v Value, ttl time.Duration) error
}
