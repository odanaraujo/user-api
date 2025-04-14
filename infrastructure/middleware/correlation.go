package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"go.uber.org/zap"
)

const correlationKey = "correlation_id"

func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := c.GetHeader("X-Correlation-ID")
		if cid == "" {
			cid = uuid.New().String()
		}

		c.Set("correlation_id", cid)

		// Use GetLogger para acessar o logger base
		baseLogger := loggers.GetLogger()
		requestLogger := baseLogger.With(zap.String("correlation_id", cid))

		ctxWithLogger := loggers.WithLogger(c.Request.Context(), requestLogger)
		c.Request = c.Request.WithContext(ctxWithLogger)

		c.Next()
	}
}
