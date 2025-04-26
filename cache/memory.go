package cache

import (
	"context"
	"strconv"
	"sync"
	"time"
)

type MemoryCache struct {
	Data        map[string][]byte
	Expiration  map[string]time.Time
	Mutex       sync.Mutex
	DefaultTTL  time.Duration
	CleanupTick time.Duration
}

// força o GO a garantir no build time que MemoryCache implementa cache. Se algum método faltar, dará erro
var _ Cache = (*MemoryCache)(nil)

func NewMemoryCache(defaultTTL, cleanupTick time.Duration) *MemoryCache {
	cache := &MemoryCache{
		Data:        make(map[string][]byte),
		Expiration:  make(map[string]time.Time),
		DefaultTTL:  defaultTTL,
		CleanupTick: cleanupTick,
	}

	go cache.startCleanUP()
	return cache
}

func (c *MemoryCache) startCleanUP() {
	ticker := time.NewTicker(c.CleanupTick)
	for range ticker.C {
		c.cleanup()
	}
}

func (c *MemoryCache) cleanup() {
	currentTime := time.Now()

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for key, expirationTime := range c.Expiration {
		if currentTime.After(expirationTime) {
			delete(c.Data, key)
			delete(c.Expiration, key)
		}
	}
}

func (cache *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	if ttl <= 0 {
		ttl = cache.DefaultTTL
	}

	cache.Data[key] = value
	cache.Expiration[key] = time.Now().Add(ttl)
}

func (cache *MemoryCache) Get(ctx context.Context, key string) ([]byte, bool) {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	exp, found := cache.Expiration[key]
	if !found || time.Now().After(exp) {
		delete(cache.Data, key)
		delete(cache.Expiration, key)
		return nil, false
	}

	val, ok := cache.Data[key]
	return val, ok
}

func (cache *MemoryCache) Delete(ctx context.Context, key string) {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	delete(cache.Data, key)
	delete(cache.Expiration, key)
}

func (m *MemoryCache) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	var count int64 = 1

	if val, ok := m.Data[key]; ok {
		parsed, err := strconv.ParseInt(string(val), 10, 64)
		if err != nil {
			return 0, err
		}
		count = parsed + 1
	}

	m.Data[key] = []byte(strconv.FormatInt(count, 10))
	m.Expiration[key] = time.Now().Add(expiration)

	return count, nil
}
