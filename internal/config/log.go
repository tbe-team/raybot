package config

import (
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/pkg/log"
)

type Log struct {
	Format    log.Format `yaml:"format"`
	Level     slog.Level `yaml:"level"`
	AddSource bool       `yaml:"add_source"`
}

func (l Log) Validate() error {
	if l.Format != log.FormatJSON && l.Format != log.FormatText {
		return fmt.Errorf("invalid format: %s", l.Format)
	}

	if l.Level < slog.LevelDebug || l.Level > slog.LevelError {
		return fmt.Errorf("invalid level: %s", l.Level)
	}

	return nil
}
