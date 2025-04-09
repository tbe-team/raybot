package config

import (
	"context"

	"github.com/tbe-team/raybot/internal/config"
)

type Service interface {
	GetLogConfig(ctx context.Context) (config.Log, error)
	UpdateLogConfig(ctx context.Context, logCfg config.Log) (config.Log, error)

	GetHardwareConfig(ctx context.Context) (config.Hardware, error)
	UpdateHardwareConfig(ctx context.Context, hardwareCfg config.Hardware) (config.Hardware, error)

	GetCloudConfig(ctx context.Context) (config.Cloud, error)
	UpdateCloudConfig(ctx context.Context, cloudCfg config.Cloud) (config.Cloud, error)

	GetGRPCConfig(ctx context.Context) (config.GRPC, error)
	UpdateGRPCConfig(ctx context.Context, grpcCfg config.GRPC) (config.GRPC, error)

	GetHTTPConfig(ctx context.Context) (config.HTTP, error)
	UpdateHTTPConfig(ctx context.Context, httpCfg config.HTTP) (config.HTTP, error)
}
