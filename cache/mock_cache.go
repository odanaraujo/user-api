package cache

import (
	"context"
	"time"
)

type mockCache struct {
	store map[string][]byte
}

func NewMockCache() *mockCache {
	return &mockCache{store: make(map[string][]byte)}
}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, bool) {
	val, ok := m.store[key]
	return val, ok
}

func (m *mockCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) {
	m.store[key] = value
}

func (m *mockCache) Delete(ctx context.Context, key string) {
	delete(m.store, key)
}

func (m *mockCache) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	return 0, nil
}
