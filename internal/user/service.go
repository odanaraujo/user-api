package user

import (
	"context"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/usecase/usercase/createcases"
	"github.com/odanaraujo/user-api/usecase/usercase/deletecases"
	"github.com/odanaraujo/user-api/usecase/usercase/getcases"
	"github.com/odanaraujo/user-api/usecase/usercase/updatecases"

	"github.com/odanaraujo/user-api/internal/model"
)

type Service interface {
	GetUserByID(ctx context.Context, id string) (*model.User, *exception.Exception)
	CreateUser(ctx context.Context, user *model.User) (*model.User, *exception.Exception)
	UpdateUser(ctx context.Context, user *model.User) *exception.Exception
	DeleteUser(ctx context.Context, id string) *exception.Exception
}

type UserService struct {
	getcases       *getcases.GetUserUseCase
	createUserCase *createcases.CreateUserUseCase
	updateUserCase *updatecases.UpdateUserUseCase
	deleteUserCase *deletecases.DeleteUserUseCase
}

func NewUserService(cache cache.Cache) *UserService {
	return &UserService{
		getcases:       getcases.NewGetUserCase(cache),
		createUserCase: createcases.NewCreateUser(cache),
		updateUserCase: updatecases.NewUpdateUser(cache),
		deleteUserCase: deletecases.NewDeleteUser(cache),
	}
}

func (u *UserService) GetUserByID(ctx context.Context, ID string) (*model.User, *exception.Exception) {
	return u.getcases.Execute(ctx, ID)
}

func (u *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, *exception.Exception) {
	return u.createUserCase.Execute(ctx, user)
}

func (u *UserService) UpdateUser(ctx context.Context, user *model.User) *exception.Exception {
	return u.updateUserCase.Execute(ctx, user)
}

func (u *UserService) DeleteUser(ctx context.Context, id string) *exception.Exception {
	return u.deleteUserCase.Execute(ctx, id)
}
