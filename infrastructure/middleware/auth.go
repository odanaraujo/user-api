package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/internal/auth"
)

func AuthMiddleware(authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, exception.UnauthorizedRequestException("authorization header is required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, exception.UnauthorizedRequestException("invalid authorization header format"))
			return
		}

		claims, err := authService.ValidateToken(c.Request.Context(), parts[1])
		if err != nil {
			c.AbortWithStatusJSON(err.Code, err)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
