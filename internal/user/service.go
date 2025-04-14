package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"go.uber.org/zap"

	"github.com/odanaraujo/user-api/internal/model"
)

const (
	cacheTTL = time.Hour * 24
)

type Service interface {
	GetUserByID(ctx context.Context, id string) (*model.User, *exception.Exception)
	CreateUser(ctx context.Context, user *model.User) (*model.User, *exception.Exception)
	UpdateUser(ctx context.Context, user *model.User) *exception.Exception
	DeleteUser(ctx context.Context, id string) *exception.Exception
}

type UserService struct {
	Cache cache.Cache
}

func NewUserService(cache cache.Cache) *UserService {
	return &UserService{Cache: cache}
}

func (u *UserService) GetUserByID(ctx context.Context, id string) (*model.User, *exception.Exception) {
	if id == "" {
		return nil, exception.BadRequestException("user ID is required")
	}

	var user model.User
	userCache, ok := u.Cache.Get(ctx, id)
	if !ok {
		log := loggers.FromContext(ctx)
		log.Error("user not found in cache", zap.String("user_id", id))
		return nil, exception.NotFoundException("user not found in cache")
	}

	if err := json.Unmarshal(userCache, &user); err != nil {
		log := loggers.FromContext(ctx)
		log.Error("error performing unmarshal", zap.Error(err), zap.String("user_id", id))
		return nil, exception.InternalServerException("error performing unmarshal")
	}

	return &user, nil
}

func (u *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, *exception.Exception) {
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

	u.Cache.Set(ctx, user.ID, data, cacheTTL)

	return user, nil

}

func (u *UserService) UpdateUser(ctx context.Context, user *model.User) *exception.Exception {
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

	u.Cache.Set(ctx, userForUpdate.ID, updatedData, cacheTTL)

	return nil

}

func (u *UserService) DeleteUser(ctx context.Context, id string) *exception.Exception {
	if id == "" {
		return exception.BadRequestException("user ID is required")
	}

	if _, ok := u.Cache.Get(ctx, id); !ok {
		return exception.NotFoundException("user not found in cache")
	}

	u.Cache.Delete(ctx, id)
	return nil
}

func (u *UserService) buildUserUpdate(newUser, oldUser *model.User) *model.User {
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
