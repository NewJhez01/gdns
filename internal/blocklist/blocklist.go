package blocklist

import "context"

type Blocklist interface {
	IsBlocked(key string, ctx context.Context) (bool, error)
}
