package updatecases

import (
	"context"
	"encoding/json"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"github.com/odanaraujo/user-api/internal/model"
	"go.uber.org/zap"
)

type UpdateUserUseCase struct {
	Cache cache.Cache
}

func NewUpdateUser(cache cache.Cache) *UpdateUserUseCase {
	return &UpdateUserUseCase{Cache: cache}
}

func (u *UpdateUserUseCase) Execute(ctx context.Context, user *model.User) *exception.Exception {
	if err := user.Validate(); err != nil {
		log := loggers.FromContext(ctx)
		log.Error("invalid user data", zap.String("user_id", user.ID))
		return exception.BadRequestException("invalid user data")
	}

	data, ok := u.Cache.Get(ctx, user.ID)
	if !ok {
		return exception.NotFoundException("user not found in cache")
	}

	var oldUser model.User
	if err := json.Unmarshal(data, &oldUser); err != nil {
		return exception.InternalServerException("user not found in cache")
	}

	userForUpdate := u.buildUserUpdate(user, &oldUser)
	updatedData, err := json.Marshal(userForUpdate)
	if err != nil {
		return exception.InternalServerException("failed to marshal updated user")
	}

	u.Cache.Set(ctx, userForUpdate.ID, updatedData, model.CacheTTL)

	return nil
}

func (u *UpdateUserUseCase) buildUserUpdate(newUser, oldUser *model.User) *model.User {
	if newUser.Name != "" {
		oldUser.Name = newUser.Name
	}

	if newUser.CPF != "" {
		oldUser.CPF = newUser.CPF
	}

	if newUser.Age != 0 {
		oldUser.Age = newUser.Age
	}

	if newUser.Email != "" {
		oldUser.Email = newUser.Email
	}

	return oldUser
}
