package logging

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/tbe-team/raybot/internal/config"
)

type CleanupFunc func() error

// NewSlogLogger creates and configures a new slog.Logger.
func NewSlogLogger(cfg config.Log) (*slog.Logger, CleanupFunc, error) {
	handlers := []slog.Handler{}
	cleanupFunc := func() error { return nil }

	if cfg.File.Enable {
		fileHandler, cleanup, err := newFileHandler(cfg.File)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create file handler: %w", err)
		}
		handlers = append(handlers, fileHandler)
		cleanupFunc = cleanup
	}

	if cfg.Console.Enable {
		consoleHandler := newConsoleHandler(cfg.Console)
		handlers = append(handlers, consoleHandler)
	}

	// Fallback no operation logger if no handlers are enabled
	if len(handlers) == 0 {
		return NewNoopLogger(), func() error { return nil }, nil
	}

	// Use Fanout if multiple handlers
	var handler slog.Handler
	if len(handlers) == 1 {
		handler = handlers[0]
	} else {
		handler = Fanout(handlers...)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, cleanupFunc, nil
}

func newConsoleHandler(cfg config.LogConsoleHandler) slog.Handler {
	if cfg.Format == config.LogFormatJSON {
		return slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: cfg.Level,
		})
	}

	return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.Level,
	})
}

func newFileHandler(cfg config.LogFileHandler) (slog.Handler, CleanupFunc, error) {
	l := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxBackups: cfg.RotationCount,
		// 1 Terabyte -- basically infinite. Don't rollover on size. Just restarts.
		MaxSize: 1024 * 1024,
	}

	// Rotate on restart
	if err := l.Rotate(); err != nil {
		return nil, nil, fmt.Errorf("failed to rotate log file: %w", err)
	}

	cleanup := func() error {
		return l.Close()
	}

	if cfg.Format == config.LogFormatJSON {
		return slog.NewJSONHandler(l, &slog.HandlerOptions{
			Level: cfg.Level,
		}), cleanup, nil
	}

	return slog.NewTextHandler(l, &slog.HandlerOptions{
		Level: cfg.Level,
	}), cleanup, nil
}

func NewNoopLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
