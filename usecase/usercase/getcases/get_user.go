package getcases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"github.com/odanaraujo/user-api/internal/model"
	"go.uber.org/zap"
)

type GetUserUseCase struct {
	Cache cache.Cache
}

func NewGetUserCase(cache cache.Cache) *GetUserUseCase {
	return &GetUserUseCase{
		Cache: cache,
	}
}

func (g *GetUserUseCase) Execute(ctx context.Context, ID string) (*model.User, *exception.Exception) {
	fmt.Printf("init get [DAN]: %s", ID)
	if ID == "" {
		return nil, exception.BadRequestException("user ID is required")
	}

	var user model.User
	userCache, ok := g.Cache.Get(ctx, ID)
	if !ok {
		log := loggers.FromContext(ctx)
		log.Error("user not found in cache", zap.String("user_id", ID))
		return nil, exception.NotFoundException("user not found in cache")
	}

	if err := json.Unmarshal(userCache, &user); err != nil {
		log := loggers.FromContext(ctx)
		log.Error("error performing unmarshal", zap.Error(err), zap.String("user_id", ID))
		return nil, exception.InternalServerException("error performing unmarshal")
	}

	return &user, nil
}
