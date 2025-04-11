package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/configurations/exception"
	"github.com/odanaraujo/user-api/internal/model"
	"github.com/odanaraujo/user-api/internal/user"
)

type UserHandler struct {
	Service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{Service: service}
}

func (u *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := u.Service.GetUserByID(id)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		ex := exception.BadRequestException("invalid request payload")
		c.JSON(ex.Code, ex)
		return
	}

	createUser, err := u.Service.CreateUser(&user)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusCreated, createUser)
}

func (u *UserHandler) UpdateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		ex := exception.BadRequestException("invalid request payload")
		c.JSON(ex.Code, ex)
		return
	}

	if ex := u.Service.UpdateUser(&user); ex != nil {
		c.JSON(ex.Code, ex)
		return
	}

	c.Status(http.StatusNoContent)
}

func (u *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := u.Service.DeleteUser(id); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusNoContent)
}
