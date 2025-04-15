package router

import (
	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/infrastructure/middleware"
	"github.com/odanaraujo/user-api/routes"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CorrelationIDMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	routes.RegisterRoutes(r)
	return r
}
