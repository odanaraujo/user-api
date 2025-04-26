package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/internal/handler"
	"github.com/odanaraujo/user-api/internal/user"
)

func RegisterRoutes(r *gin.Engine, userService user.Service) {
	gin.Default()

	userHandler := handler.NewUserHandler(userService)

	userGroup := r.Group("/users")

	{
		userGroup.GET("/:id", userHandler.GetUserByID)
		userGroup.POST("", userHandler.CreateUser)
		userGroup.PUT("", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
	}

}
