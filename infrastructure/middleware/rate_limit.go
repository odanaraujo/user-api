package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"golang.org/x/time/rate"
)

var (
	limiter *rate.Limiter
	mu      sync.Mutex
)

func init() {
	limiter = rate.NewLimiter(rate.Every(time.Second/10), 5)
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		// check if the requested quantity is within the limit
		if !limiter.Allow() {
			err := exception.TooManyRequestsException("Too many requests, please try again later.")
			log := loggers.FromContext(c.Request.Context())
			log.Error(err.Message)
			c.JSON(err.Code, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
