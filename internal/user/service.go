package user

import (
	"encoding/json"
	"time"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/configurations/exception"

	"github.com/odanaraujo/user-api/internal/model"
)

const (
	cacheTTL = time.Hour * 24
)

type Service interface {
	GetUserByID(id string) (*model.User, *exception.Exception)
	CreateUser(user *model.User) (*model.User, *exception.Exception)
	UpdateUser(user *model.User) *exception.Exception
	DeleteUser(id string) *exception.Exception
}

type UserService struct {
	Cache cache.Cache
}

func NewUserService(cache cache.Cache) *UserService {
	return &UserService{Cache: cache}
}

func (u *UserService) GetUserByID(id string) (*model.User, *exception.Exception) {
	if id == "" {
		return nil, exception.BadRequestException("user ID is required")
	}

	var user model.User
	userCache, ok := u.Cache.Get(id)
	if !ok {
		return nil, exception.NotFoundException("user not found in cache")
	}

	if err := json.Unmarshal(userCache, &user); err != nil {
		return nil, exception.InternalServerException("error performing unmarshal")
	}

	return &user, nil
}

func (u *UserService) CreateUser(user *model.User) (*model.User, *exception.Exception) {
	if err := user.Validate(); err != nil {
		return nil, exception.BadRequestException("invalid user data")
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, exception.InternalServerException("could not persist user in cache")
	}

	u.Cache.Set(user.ID, data, cacheTTL)

	return user, nil

}

func (u *UserService) UpdateUser(user *model.User) *exception.Exception {
	if err := user.Validate(); err != nil {
		return exception.BadRequestException("invalid user data")
	}

	data, ok := u.Cache.Get(user.ID)
	if !ok {
		return exception.NotFoundException("user not found in cache")
	}

	var oldUser model.User
	if err := json.Unmarshal(data, &oldUser); err != nil {
		return exception.InternalServerException("failed to unmarshal user with ID")
	}

	userForUpdate := u.buildUserUpdate(user, &oldUser)
	updatedData, err := json.Marshal(userForUpdate)
	if err != nil {
		return exception.InternalServerException("failed to marshal updated user")
	}

	u.Cache.Set(userForUpdate.ID, updatedData, cacheTTL)

	return nil

}

func (u *UserService) DeleteUser(id string) *exception.Exception {
	if id == "" {
		return exception.BadRequestException("user ID is required")
	}

	if _, ok := u.Cache.Get(id); !ok {
		return exception.NotFoundException("user not found in cache")
	}

	u.Cache.Delete(id)
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
