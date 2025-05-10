package usercase

import (
	"github.com/odanaraujo/user-api/cache"
)

type ServiceImpl struct {
	Cache cache.Cache
}

func NewUserService(cache cache.Cache) *ServiceImpl {
	return &ServiceImpl{Cache: cache}
}
