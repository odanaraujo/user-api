package deletecases

import (
	"context"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/exception"
)

type DeleteUserUseCase struct {
	Cache cache.Cache
}

func NewDeleteUser(cache cache.Cache) *DeleteUserUseCase {
	return &DeleteUserUseCase{Cache: cache}
}

func (d *DeleteUserUseCase) Execute(ctx context.Context, id string) *exception.Exception {
	if id == "" {
		return exception.BadRequestException("user ID is required")
	}

	if _, ok := d.Cache.Get(ctx, id); !ok {
		return exception.NotFoundException("user not found in cache")
	}

	d.Cache.Delete(ctx, id)
	return nil
}
