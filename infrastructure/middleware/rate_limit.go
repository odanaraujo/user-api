package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"go.uber.org/zap"
)

const (
	rateLimitMaxRequests = 5
	rateLimitWindow      = time.Minute * 1
)

func RateLimitByIP(cache cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		count, err := cache.Increment(ctx, key, rateLimitWindow)
		if err != nil {
			loggers.FromContext(ctx).Error("error incrementing rate limit", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, exception.InternalServerException("internal error"))
			return
		}

		if count > int64(rateLimitMaxRequests) {
			loggers.FromContext(ctx).Warn("rate limit exceeded", zap.String("ip", ip))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, exception.TooManyRequestsException("Too many requests, try again later"))
			return
		}

		c.Next()
	}
}
