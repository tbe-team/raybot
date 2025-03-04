package log

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type loggerCtxKey int

// CtxLoggerKey is the context key for the logger.
const CtxLoggerKey loggerCtxKey = iota

// WithLogger returns a new context that includes the provided logger.
// Useful for propagating logger configuration across different parts of the application.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, CtxLoggerKey, logger)
}

// FromContext returns the logger from context.
// If no logger is found, it returns the default logger.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(CtxLoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// NewLogger creates a new logger based on the provided configuration.
func NewLogger(cfg Config) *slog.Logger {
	var handler slog.Handler

	level := slog.LevelInfo
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:     level,
			AddSource: cfg.AddSource,
		})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     level,
			AddSource: cfg.AddSource,
		})
	}

	return slog.New(handler)
}

// CloneLogger creates a new logger with the same configuration as the provided logger.
func CloneLogger(logger *slog.Logger) *slog.Logger {
	return slog.New(logger.Handler())
}
