package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/infrastructure/middleware"
	"github.com/odanaraujo/user-api/internal/handler"
	"github.com/odanaraujo/user-api/internal/user"
)

func RegisterRoutes(r *gin.Engine, userService user.Service) {
	gin.Default()

	userHandler := handler.NewUserHandler(userService)

	userGroup := r.Group("/users")
	{
		// Public routes
		userGroup.POST("", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUserByID)

		// Protected routes
		protected := userGroup.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.PUT("", userHandler.UpdateUser)
			protected.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}
