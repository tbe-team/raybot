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
		return config.Log{}, fmt.Errorf("validate log config: %w", err)
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
		return config.Hardware{}, fmt.Errorf("validate hardware config: %w", err)
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
		return config.Cloud{}, fmt.Errorf("validate cloud config: %w", err)
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

func (s *service) GetGRPCConfig(_ context.Context) (config.GRPC, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.GRPC, nil
}

func (s *service) UpdateGRPCConfig(ctx context.Context, grpcCfg config.GRPC) (config.GRPC, error) {
	if err := grpcCfg.Validate(); err != nil {
		return config.GRPC{}, fmt.Errorf("validate grpc config: %w", err)
	}

	cfg := *s.cfg
	cfg.GRPC = grpcCfg

	if err := s.writeConfig(ctx, cfg); err != nil {
		return config.GRPC{}, fmt.Errorf("write config: %w", err)
	}

	s.mu.Lock()
	s.cfg = &cfg
	s.mu.Unlock()

	return grpcCfg, nil
}

func (s *service) GetHTTPConfig(_ context.Context) (config.HTTP, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.HTTP, nil
}

func (s *service) UpdateHTTPConfig(ctx context.Context, httpCfg config.HTTP) (config.HTTP, error) {
	if err := httpCfg.Validate(); err != nil {
		return config.HTTP{}, fmt.Errorf("validate http config: %w", err)
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
