package service

import (
	"context"
	"time"
)

type SerialConfig struct {
	Port        string
	BaudRate    int
	DataBits    int
	StopBits    float64
	Parity      string
	ReadTimeout time.Duration
}

type PICConfig struct {
	Serial SerialConfig
}

type GRPCConfig struct {
	Port int
}

type HTTPConfig struct {
	EnableSwagger bool
}

type LogConfig struct {
	Level     string
	Format    string
	AddSource bool
}

type GetSystemConfigOutput struct {
	LogConfig  LogConfig
	PICConfig  PICConfig
	GRPCConfig GRPCConfig
	HTTPConfig HTTPConfig
}

type UpdateSystemConfigParams struct {
	LogConfig  LogConfig
	PICConfig  PICConfig
	GRPCConfig GRPCConfig
	HTTPConfig HTTPConfig
}

type UpdateSystemConfigOutput = GetSystemConfigOutput

type SystemService interface {
	GetSystemConfig(ctx context.Context) (GetSystemConfigOutput, error)
	UpdateSystemConfig(ctx context.Context, params UpdateSystemConfigParams) (UpdateSystemConfigOutput, error)

	// RestartApplication restarts the application after 3 second.
	RestartApplication(ctx context.Context) error
}
