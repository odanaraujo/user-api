package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisInstance *redisCache
	once          sync.Once
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache() *redisCache {
	once.Do(func() {
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		password := os.Getenv("REDIS_PASSWORD")
		dbStr := os.Getenv("REDIS_DB")

		db, err := strconv.Atoi(dbStr)
		if err != nil {
			db = 0
		}

		client := redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: password,
			DB:       db,
		})

		redisInstance = &redisCache{client: client}
	})

	return redisInstance
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, bool) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	return val, true
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) {
	cmd := r.client.Set(ctx, key, value, ttl)

	if cmd.Err() != nil {
		fmt.Printf("error in set redis: %v", cmd.Err().Error())
	}
}

func (r *redisCache) Delete(ctx context.Context, key string) {
	r.client.Del(ctx, key)
}
