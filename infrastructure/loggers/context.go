package loggers

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	correlationKey = "correlation_id"
	loggerKey      = contextKey("logger")
)

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extrai o correlation_id e retorna um logger com ele embutido
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return log
	}

	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok && l != nil {
		return l
	}

	return log
}

// Para uso com context.Context tradicional (caso use fora do Gin)
func FromStdContext(ctx context.Context) *zap.Logger {
	if cid, ok := ctx.Value(correlationKey).(string); ok && cid != "" {
		return log.With(zap.String(correlationKey, cid))
	}
	return log
}
