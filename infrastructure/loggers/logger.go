package loggers

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

const (
	JSON_TYPE   = "json"
	INFO_LEVEL  = "info"
	ERROR_LEVEL = "error"
	DEBUG_LEVEL = "debug"
	LOG_OUTPUT  = "LOG_OUTPUT"
	LOG_LEVEL   = "LOG_LEVEL"
)

func init() {

	logConfig := zap.Config{
		OutputPaths: []string{getOutputLogs()},
		Level:       zap.NewAtomicLevelAt(getLevelLogs()),
		Encoding:    JSON_TYPE,
		EncoderConfig: zapcore.EncoderConfig{
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder, //padroniza os loggers de alguma forma.
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	log, _ = logConfig.Build()

}

func Info(message string, tags ...zap.Field) {
	log.Info(message, tags...)
}

func Error(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Error(message, tags...)
}

// get the environment configuration
func getOutputLogs() string {
	output := strings.ToLower(strings.TrimSpace(os.Getenv(LOG_OUTPUT)))

	if output == "" {
		return "stdout"
	}

	return output
}

// verify and set the level
func getLevelLogs() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(LOG_LEVEL))) {
	case INFO_LEVEL:
		return zapcore.InfoLevel
	case ERROR_LEVEL:
		return zapcore.ErrorLevel
	case DEBUG_LEVEL:
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func Close() {
	log.Sync()
}
