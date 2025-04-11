package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/internal/handler"
	"github.com/odanaraujo/user-api/internal/user"
)

const (
	defaultTTL      = time.Hour * 24
	cleanupInterval = 10 * time.Minute
)

func RegisterRoutes(r *gin.Engine) {
	gin.Default()
	memRouter := cache.NewMemoryCache(defaultTTL, cleanupInterval)
	userService := user.NewUserService(memRouter)
	userHandler := handler.NewUserHandler(userService)

	userGroup := r.Group("/users")

	{
		userGroup.GET("/:id", userHandler.GetUserByID)
		userGroup.POST("", userHandler.CreateUser)
		userGroup.PUT("", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
	}

}
