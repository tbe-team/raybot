package configimpl

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/tbe-team/raybot/internal/config"
	configsvc "github.com/tbe-team/raybot/internal/services/config"
	"github.com/tbe-team/raybot/internal/storage/file"
	"github.com/tbe-team/raybot/pkg/xerror"
)

type service struct {
	cfg *config.Config
	mu  sync.RWMutex

	fileClient file.Client
}

func NewService(cfg *config.Config, fileClient file.Client) configsvc.Service {
	return &service{
		cfg:        cfg,
		fileClient: fileClient,
	}
}

func (s *service) GetLogConfig(_ context.Context) (config.Log, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Log, nil
}

func (s *service) UpdateLogConfig(ctx context.Context, logCfg config.Log) (config.Log, error) {
	if err := logCfg.Validate(); err != nil {
		return config.Log{}, xerror.ValidationFailed(err, "invalid log config")
	}

	cfg := *s.cfg
	cfg.Log = logCfg

	if err := s.writeConfig(ctx, cfg); err != nil {
		return config.Log{}, fmt.Errorf("write config: %w", err)
	}

	s.mu.Lock()
	s.cfg = &cfg
	s.mu.Unlock()

	return logCfg, nil
}

func (s *service) GetHardwareConfig(_ context.Context) (config.Hardware, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Hardware, nil
}

func (s *service) UpdateHardwareConfig(ctx context.Context, hardwareCfg config.Hardware) (config.Hardware, error) {
	if err := hardwareCfg.Validate(); err != nil {
		return config.Hardware{}, xerror.ValidationFailed(err, "invalid hardware config")
	}

	cfg := *s.cfg
	cfg.Hardware = hardwareCfg

	if err := s.writeConfig(ctx, cfg); err != nil {
		return config.Hardware{}, fmt.Errorf("write config: %w", err)
	}

	s.mu.Lock()
	s.cfg = &cfg
	s.mu.Unlock()

	return hardwareCfg, nil
}

func (s *service) GetCloudConfig(_ context.Context) (config.Cloud, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Cloud, nil
}

func (s *service) UpdateCloudConfig(ctx context.Context, cloudCfg config.Cloud) (config.Cloud, error) {
	if err := cloudCfg.Validate(); err != nil {
		return config.Cloud{}, xerror.ValidationFailed(err, "invalid cloud config")
	}

	cfg := *s.cfg
	cfg.Cloud = cloudCfg

	if err := s.writeConfig(ctx, cfg); err != nil {
		return config.Cloud{}, fmt.Errorf("write config: %w", err)
	}

	s.mu.Lock()
	s.cfg = &cfg
	s.mu.Unlock()

	return cloudCfg, nil
}

func (s *service) GetHTTPConfig(_ context.Context) (config.HTTP, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.HTTP, nil
}

func (s *service) UpdateHTTPConfig(ctx context.Context, httpCfg config.HTTP) (config.HTTP, error) {
	if err := httpCfg.Validate(); err != nil {
		return config.HTTP{}, xerror.ValidationFailed(err, "invalid http config")
	}

	cfg := *s.cfg
	cfg.HTTP = httpCfg

	if err := s.writeConfig(ctx, cfg); err != nil {
		return config.HTTP{}, fmt.Errorf("write config: %w", err)
	}

	s.mu.Lock()
	s.cfg = &cfg
	s.mu.Unlock()

	return httpCfg, nil
}

func (s *service) GetWifiConfig(_ context.Context) (config.Wifi, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Wifi, nil
}

func (s *service) UpdateWifiConfig(ctx context.Context, wifiCfg config.Wifi) (config.Wifi, error) {
	if err := wifiCfg.Validate(); err != nil {
		return config.Wifi{}, xerror.ValidationFailed(err, "invalid wifi config")
	}

	cfg := *s.cfg
	cfg.Wifi = wifiCfg

	if err := s.writeConfig(ctx, cfg); err != nil {
		return config.Wifi{}, fmt.Errorf("write config: %w", err)
	}

	s.mu.Lock()
	s.cfg = &cfg
	s.mu.Unlock()

	return wifiCfg, nil
}

func (s *service) GetCommandConfig(_ context.Context) (config.Command, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Command, nil
}

func (s *service) UpdateCommandConfig(ctx context.Context, commandCfg config.Command) (config.Command, error) {
	if err := commandCfg.Validate(); err != nil {
		return config.Command{}, xerror.ValidationFailed(err, "invalid command config")
	}

	cfg := *s.cfg
	cfg.Command = commandCfg

	if err := s.writeConfig(ctx, cfg); err != nil {
		return config.Command{}, fmt.Errorf("write config: %w", err)
	}

	s.mu.Lock()
	s.cfg = &cfg
	s.mu.Unlock()

	return commandCfg, nil
}

func (s *service) writeConfig(ctx context.Context, cfg config.Config) error {
	writer, err := s.fileClient.Write(ctx, s.cfg.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	defer writer.Close()

	buf := bytes.Buffer{}
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	if _, err := writer.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
