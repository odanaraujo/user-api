package router

import (
	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/middleware"
	"github.com/odanaraujo/user-api/internal/user"
	"github.com/odanaraujo/user-api/routes"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CorrelationIDMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// inject user service
	redis := cache.NewRedisCache()
	userService := user.NewUserService(redis)

	routes.RegisterRoutes(r, userService)
	return r
}
