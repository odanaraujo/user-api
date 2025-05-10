package router

import (
	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/middleware"
	"github.com/odanaraujo/user-api/internal/auth"
	"github.com/odanaraujo/user-api/internal/user"
	"github.com/odanaraujo/user-api/routes"
)

func NewRouter() *gin.Engine {
	// inject dependencies
	redis := cache.NewRedisCache()
	authService := auth.NewAuthService(redis)
	userService := user.NewUserService(redis, authService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CorrelationIDMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RateLimitByIP(redis))

	routes.RegisterRoutes(r, userService, authService)
	return r
}
