package usercase

import (
	"context"

	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/internal/model"
)

type UseCase interface {
	GetUserByID(ctx context.Context, id string) (*model.UserResponse, *exception.Exception)
	CreateUser(ctx context.Context, user *model.User) (*model.User, *exception.Exception)
	UpdateUser(ctx context.Context, user *model.User) *exception.Exception
	DeleteUser(ctx context.Context, id string) *exception.Exception
}
