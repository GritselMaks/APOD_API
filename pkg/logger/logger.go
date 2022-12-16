package logger

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Private type for context map
type contextKey string

const (
	levelDebug     = "DEBUG"
	levelInfo      = "INFO"
	levelWarning   = "WARNING"
	levelError     = "ERROR"
	levelCritical  = "CRITICAL"
	levelAlert     = "ALERT"
	levelEmergency = "EMERGENCY"

	encodingConsole = "console"
	encodingJSON    = "json"

	// loggerKey points to the value in the context where the logger is stored.
	loggerKey = contextKey("logger")
)

var outputStderr = []string{"stderr"}

type Logger struct {
	*zap.SugaredLogger
}

// Returns reasonable development logging configuration.
// NewDevelopmentEncoderConfig returns an opinionated EncoderConfig for
// development environments.
func NewLogger(level string) *Logger {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(levelToZapLevel(level)),
		Development:      true,
		Encoding:         encodingConsole,
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      outputStderr,
		ErrorOutputPaths: outputStderr,
	}
	l, err := config.Build()
	if err != nil {
		l = zap.NewNop()
	}
	return &Logger{l.Sugar()}
}

// levelToZapLevel converts the given string to the appropriate zap level value.
func levelToZapLevel(s string) zapcore.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case levelDebug:
		return zapcore.DebugLevel
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarning:
		return zapcore.WarnLevel
	case levelError:
		return zapcore.ErrorLevel
	case levelCritical:
		return zapcore.DPanicLevel
	case levelAlert:
		return zapcore.PanicLevel
	case levelEmergency:
		return zapcore.FatalLevel
	}
	return zapcore.WarnLevel
}

// CtxWithLogger creates a new context with the provided logger attached.
func CtxWithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Return logger stored in context
func LoggerFromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return &Logger{logger}
	}
	return NewLogger("")
}
