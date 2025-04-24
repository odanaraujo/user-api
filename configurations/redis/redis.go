package redis

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	clientInstance *redis.Client
	once           sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		clientInstance = redis.NewClient(&redis.Options{
			Addr:         getRedisAddr(),
			Password:     getRedisPassword(),
			DB:           0,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})

		if err := clientInstance.Ping(context.Background()).Err(); err != nil {
			panic(fmt.Sprintf("error connecting to redis: %v", err))
		}
	})
	return clientInstance
}

func getRedisAddr() string {
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		return addr
	}

	return "localhost:6379"
}

func getRedisPassword() string {
	if pwd := os.Getenv("REDIS_PASSWORD"); pwd != "" {
		return pwd
	}

	return ""
}
