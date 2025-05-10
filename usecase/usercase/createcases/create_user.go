package createcases

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"github.com/odanaraujo/user-api/internal/auth"
	"github.com/odanaraujo/user-api/internal/model"
	"go.uber.org/zap"
)

type CreateUserUseCase struct {
	Cache cache.Cache
	Auth  auth.Service
}

func NewCreateUser(cache cache.Cache, authService auth.Service) *CreateUserUseCase {
	return &CreateUserUseCase{
		Cache: cache,
		Auth:  authService,
	}
}

func (c *CreateUserUseCase) Execute(ctx context.Context, user *model.User) (*model.CreateUserResponse, *exception.Exception) {
	user.ID = uuid.NewString()

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

	token, ex := c.Auth.GenerateToken(ctx, user.ID)
	if ex != nil {
		return nil, ex
	}

	return &model.CreateUserResponse{
		User:  model.NewUserResponse(user),
		Token: token,
	}, nil
}
