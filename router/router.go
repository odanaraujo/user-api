package router

import (
	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/configurations/middleware"
	"github.com/odanaraujo/user-api/routes"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	routes.RegisterRoutes(r)
	return r
}
