package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/odanaraujo/user-api/infrastructure/exception"
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
	ctx := c.Request.Context()
	id := c.Param("id")
	user, err := u.Service.GetUserByID(ctx, id)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var user model.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		ex := exception.BadRequestException("invalid request payload")
		c.JSON(ex.Code, ex)
		return
	}

	response, err := u.Service.CreateUser(ctx, &user)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (u *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, exception.UnauthorizedRequestException("unauthorized"))
		return
	}

	var user model.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		ex := exception.BadRequestException("invalid request payload")
		c.JSON(ex.Code, ex)
		return
	}

	// Ensure user can only update their own data
	if user.ID != userID.(string) {
		c.JSON(http.StatusForbidden, exception.ForbiddenException("you can only update your own data"))
		return
	}

	if ex := u.Service.UpdateUser(ctx, &user); ex != nil {
		c.JSON(ex.Code, ex)
		return
	}

	c.Status(http.StatusNoContent)
}

func (u *UserHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, exception.UnauthorizedRequestException("unauthorized"))
		return
	}

	id := c.Param("id")
	// Ensure user can only delete their own data
	if id != userID.(string) {
		c.JSON(http.StatusForbidden, exception.ForbiddenException("you can only delete your own data"))
		return
	}

	if err := u.Service.DeleteUser(ctx, id); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusNoContent)
}
