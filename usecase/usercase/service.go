package usercase

import (
	"time"

	"github.com/odanaraujo/user-api/cache"
)

const (
	cacheTTL = time.Hour * 24
)

type ServiceImpl struct {
	Cache cache.Cache
}

func NewUserService(cache cache.Cache) *ServiceImpl {
	return &ServiceImpl{Cache: cache}
}
