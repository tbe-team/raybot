package log

import (
	"flag"
	"fmt"
	"strings"
)

// Config holds configuration for the logger.
type Config struct {
	Level     string `yaml:"level"`
	Format    string `yaml:"format"`
	AddSource bool   `yaml:"add_source"`
}

// RegisterFlags registers flags for the logger.
func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&c.Level, "log-level", "info", "log level")
	f.StringVar(&c.Format, "log-format", "text", "log format")
	f.BoolVar(&c.AddSource, "log-add-source", false, "add source to log")
}

// Validate validates the logger configuration.
func (c *Config) Validate() error {
	level := strings.ToLower(c.Level)
	if level != "debug" && level != "info" && level != "warn" && level != "error" {
		return fmt.Errorf("invalid log level: %s", c.Level)
	}

	format := strings.ToLower(c.Format)
	if format != "json" && format != "text" {
		return fmt.Errorf("invalid log format: %s", c.Format)
	}

	return nil
}
