package config

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tbe-team/raybot/internal/storage/file"
)

var (
	ErrInvalidConfig = errors.New("invalid config")
)

type Manager interface {
	GetConfig() Config
	SaveConfig(ctx context.Context, cfg Config) error
}

type manager struct {
	cfg        *Config
	configPath string
	fileClient file.Client
	log        *slog.Logger
}

// NewManager creates a new config manager.
// It detects the install path of the application and loads the config file.
// If the config file does not exist, it creates it and saves the default config.
// If the config file exists, it loads the config file.
func NewManager(fileClient file.Client, configPath string, log *slog.Logger) (Manager, error) {
	s := &manager{
		cfg:        &DefaultConfig,
		configPath: configPath,
		fileClient: fileClient,
		log:        log,
	}

	if _, err := os.Stat(s.configPath); os.IsNotExist(err) {
		if err := s.SaveConfig(context.Background(), *s.cfg); err != nil {
			return nil, fmt.Errorf("save default config: %w", err)
		}
	} else {
		if err := s.loadCfg(); err != nil {
			return nil, fmt.Errorf("load config: %w", err)
		}
	}

	// check if the config file is valid
	if err := s.cfg.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	return s, nil
}

// GetConfig returns the config.
func (s manager) GetConfig() Config {
	return *s.cfg
}

func (s *manager) SaveConfig(ctx context.Context, cfg Config) error {
	// validate the config
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	writer, err := s.fileClient.Write(ctx, s.configPath)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	defer writer.Close()

	buf := bytes.Buffer{}
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	if _, err := writer.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	s.cfg = &cfg

	return nil
}

func (s *manager) loadCfg() error {
	reader, err := s.fileClient.Read(context.Background(), s.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &s.cfg); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	return nil
}
