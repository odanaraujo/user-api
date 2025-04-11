package loggers

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const correlationKey = "correlation_id"

// FromContext extrai o correlation_id e retorna um logger com ele embutido
func FromContext(c *gin.Context) *zap.Logger {
	cid := c.GetString(correlationKey)
	if cid == "" {
		return log
	}

	return log.With(zap.String(correlationKey, cid))
}

// Para uso com context.Context tradicional (caso use fora do Gin)
func FromStdContext(ctx context.Context) *zap.Logger {
	if cid, ok := ctx.Value(correlationKey).(string); ok && cid != "" {
		return log.With(zap.String(correlationKey, cid))
	}
	return log
}
