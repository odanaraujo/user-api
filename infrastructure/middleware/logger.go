package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		correlationID, _ := ctx.Get(correlationKey)

		loggers.Info(
			"incoming request",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.RequestURI),
			zap.Int("status", ctx.Writer.Status()),
			zap.String("client_ip", ctx.ClientIP()),
			zap.Duration("duration", duration),
			zap.String("user_agente", ctx.Request.UserAgent()),
			zap.String("correlation_id", correlationID.(string)),
		)
	}
}
