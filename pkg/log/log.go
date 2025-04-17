package log

import (
	"log/slog"
	"os"
)

// NewSlogLogger creates a new slog logger.
func NewSlogLogger(cfg Config) *slog.Logger {
	var handler slog.Handler
	if cfg.Format == FormatJSON {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:     cfg.Level,
			AddSource: cfg.AddSource,
		})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     cfg.Level,
			AddSource: cfg.AddSource,
		})
	}

	log := slog.New(handler)
	slog.SetDefault(log)

	return log
}

func NewNoopLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
