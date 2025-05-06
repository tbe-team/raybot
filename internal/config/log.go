package config

import (
	"fmt"
	"log/slog"
	"strings"
)

type Log struct {
	Format    LogFormat  `yaml:"format"`
	Level     slog.Level `yaml:"level"`
	AddSource bool       `yaml:"add_source"`
}

func (l Log) Validate() error {
	if l.Format != LogFormatJSON && l.Format != LogFormatText {
		return fmt.Errorf("invalid format: %s", l.Format)
	}

	if l.Level < slog.LevelDebug || l.Level > slog.LevelError {
		return fmt.Errorf("invalid level: %s", l.Level)
	}

	return nil
}

// LogFormat represents the logging format (JSON or Text).
type LogFormat uint8

func (f LogFormat) String() string {
	return []string{"JSON", "TEXT"}[f]
}

const (
	LogFormatJSON LogFormat = iota
	LogFormatText
)

// UnmarshalText implements [encoding.TextUnmarshaler].
// It unmarshals the text to a log format.
func (f *LogFormat) UnmarshalText(text []byte) error {
	switch strings.ToUpper(string(text)) {
	case "JSON":
		*f = LogFormatJSON
	case "TEXT":
		*f = LogFormatText
	default:
		return fmt.Errorf("unknown log format: %s", text)
	}
	return nil
}

func (f LogFormat) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}
