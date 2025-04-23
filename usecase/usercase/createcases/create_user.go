package createcases

import (
	"context"
	"encoding/json"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"github.com/odanaraujo/user-api/internal/model"
	"go.uber.org/zap"
)

type CreateUserUseCase struct {
	Cache cache.Cache
}

func NewCreateUser(cache cache.Cache) *CreateUserUseCase {
	return &CreateUserUseCase{Cache: cache}
}

func (c *CreateUserUseCase) Execute(ctx context.Context, user *model.User) (*model.User, *exception.Exception) {
	if err := user.Validate(); err != nil {
		log := loggers.FromContext(ctx)
		log.Error("invalid user data", zap.String("user_id", user.ID))
		return nil, exception.BadRequestException("invalid user data")
	}

	data, err := json.Marshal(user)
	if err != nil {
		log := loggers.FromContext(ctx)
		log.Error("could not persist user in cache", zap.String("user_id", user.ID))
		return nil, exception.InternalServerException("could not persist user in cache")
	}

	c.Cache.Set(ctx, user.ID, data, model.CacheTTL)

	return user, nil
}
