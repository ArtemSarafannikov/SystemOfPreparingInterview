package logger

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/tirion/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxLoggerKey struct{}
type ctxHandlerKey struct{}

// Logger is a global zap.Logger.
var Logger *zap.Logger

// Setup initializes the global logger based on environment.
func Setup(env config.Environment) error {
	var cfg zap.Config

	switch env {
	case config.Local:
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case config.Stage:
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case config.Production:
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	default:
		return fmt.Errorf("not supported environment: %s", env)
	}

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		return fmt.Errorf("cfg.Build: %w", err)
	}

	return nil
}

// WithLogger returns a new context with the given logger stored inside it.
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, l)
}

// FromContext extracts the logger from the context. If not found, returns the global Logger.
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLoggerKey{}).(*zap.Logger); ok && l != nil {
		return l
	}
	return Logger
}

// WithHandler returns a new context with the handler name stored inside it.
// This is typically called by the gRPC interceptor automatically.
func WithHandler(ctx context.Context, handler string) context.Context {
	return context.WithValue(ctx, ctxHandlerKey{}, handler)
}

// handlerFromContext extracts the handler name from the context.
func handlerFromContext(ctx context.Context) string {
	if h, ok := ctx.Value(ctxHandlerKey{}).(string); ok {
		return h
	}
	return ""
}

// Infof logs an info message. The handler field is auto-populated from the context.
func Infof(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Info(msg, withHandler(ctx, fields)...)
}

// Warnf logs a warning message. The handler field is auto-populated from the context.
func Warnf(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Warn(msg, withHandler(ctx, fields)...)
}

// Errorf logs an error message. The handler field is auto-populated from the context.
func Errorf(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Error(msg, withHandler(ctx, fields)...)
}

// Fatalf logs a fatal message, then calls os.Exit(1). The handler field is auto-populated from the context.
func Fatalf(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Fatal(msg, withHandler(ctx, fields)...)
}

func withHandler(ctx context.Context, fields []zap.Field) []zap.Field {
	if h := handlerFromContext(ctx); h != "" {
		return append([]zap.Field{zap.String("handler", h)}, fields...)
	}
	return fields
}
