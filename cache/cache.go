package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration)
	Get(ctx context.Context, key string) ([]byte, bool)
	Delete(ctx context.Context, key string)
	Increment(ctx context.Context, key string, expiration time.Duration) (int64, error)
}
