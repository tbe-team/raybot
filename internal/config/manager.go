package config

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// configFileName is the name of the config file.
	configFileName = "config.yml"
)

var (
	ErrInvalidConfig = errors.New("invalid config")

	_ Manager = (*DefaultManager)(nil)
)

type Manager interface {
	GetConfig() Config
	SetConfig(cfg Config) error
	SaveConfig() error
	LoadConfig() error
}

type DefaultManager struct {
	cfg Config

	configPath string
	log        *slog.Logger
}

// NewManager creates a new config manager.
// It detects the install path of the application and loads the config file.
// If the config file does not exist, it creates it and saves the default config.
// If the config file exists, it loads the config file.
func NewManager() (*DefaultManager, error) {
	installPath := detectInstallPath()
	configPath := filepath.Join(installPath, configFileName)

	s := &DefaultManager{
		cfg:        DefaultConfig,
		configPath: configPath,
		log:        slog.Default(),
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		s.log.Info("Config file not found, creating it", slog.String("configPath", configPath))
		if err := initDirs(installPath); err != nil {
			return nil, fmt.Errorf("create config directory: %w", err)
		}

		if err := s.SaveConfig(); err != nil {
			return nil, fmt.Errorf("save default config: %w", err)
		}
	} else {
		if err := s.LoadConfig(); err != nil {
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
func (s *DefaultManager) GetConfig() Config {
	return s.cfg
}

// SetConfig sets the config. It does not save the config to the file.
// It is recommended to use SaveConfig() to save the config to the file.
func (s *DefaultManager) SetConfig(cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	s.cfg = cfg
	return nil
}

// LoadConfig loads the config from the file.
func (s *DefaultManager) LoadConfig() error {
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	var temp Config
	if err := yaml.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	if err := temp.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	s.cfg = temp

	return nil
}

// SaveConfig saves the config to the file.
func (s *DefaultManager) SaveConfig() error {
	// yaml.Marshal() indent 4 by default, so we use a custom encoder to indent 2
	buf := bytes.Buffer{}
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(&s.cfg); err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(s.configPath, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
