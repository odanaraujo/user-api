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
	redisInstance *RedisCache
	once          sync.Once
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
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

		redisInstance = &RedisCache{client: client}
	})

	return redisInstance
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, bool) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	return val, true
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) {
	cmd := r.client.Set(ctx, key, value, ttl)

	if cmd.Err() != nil {
		fmt.Printf("error in set redis: %v", cmd.Err().Error())
	}
}

func (r *RedisCache) Delete(ctx context.Context, key string) {
	r.client.Del(ctx, key)
}

func (r *RedisCache) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		err := r.client.Expire(ctx, key, expiration).Err()
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}
