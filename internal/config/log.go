package config

import (
	"fmt"
	"log/slog"
	"strings"
)

type Log struct {
	File    LogFileHandler    `yaml:"file"`
	Console LogConsoleHandler `yaml:"console"`
}

func (l Log) Validate() error {
	if err := l.File.Validate(); err != nil {
		return fmt.Errorf("invalid file: %w", err)
	}

	if err := l.Console.Validate(); err != nil {
		return fmt.Errorf("invalid console: %w", err)
	}

	return nil
}

const (
	defaultLogFilePath   = "logs/raybot.log"
	defaultRotationCount = 10
)

type LogFileHandler struct {
	Enable bool `yaml:"enable"`

	Path          string `yaml:"path"`
	RotationCount int    `yaml:"rotation_count"`

	Format LogFormat  `yaml:"format"`
	Level  slog.Level `yaml:"level"`
}

func (l *LogFileHandler) Validate() error {
	if l.Path == "" {
		l.Path = defaultLogFilePath
	}

	if l.RotationCount <= 0 {
		l.RotationCount = defaultRotationCount
	}

	if l.Format != LogFormatJSON && l.Format != LogFormatText {
		return fmt.Errorf("invalid format: %s", l.Format)
	}

	return nil
}

type LogConsoleHandler struct {
	Enable bool `yaml:"enable"`

	Format LogFormat  `yaml:"format"`
	Level  slog.Level `yaml:"level"`
}

func (l LogConsoleHandler) Validate() error {
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
