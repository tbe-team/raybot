package config

import (
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/controller/cloud"
	"github.com/tbe-team/raybot/internal/controller/espserial"
	"github.com/tbe-team/raybot/internal/controller/grpc"
	"github.com/tbe-team/raybot/internal/controller/http"
	"github.com/tbe-team/raybot/internal/controller/picserial"
	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
	"github.com/tbe-team/raybot/pkg/log"
)

type GRPCServerConfig = grpc.Config

type GRPCConfig struct {
	Server GRPCServerConfig `yaml:"server"`
	Cloud  cloud.Config     `yaml:"cloud"`
}

func (c *GRPCConfig) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("validate GRPC server: %w", err)
	}

	if err := c.Cloud.Validate(); err != nil {
		return fmt.Errorf("validate cloud: %w", err)
	}

	return nil
}

type Config struct {
	Log  log.Config       `yaml:"log"`
	GRPC GRPCConfig       `yaml:"grpc"`
	HTTP http.Config      `yaml:"http"`
	PIC  picserial.Config `yaml:"pic"`
	ESP  espserial.Config `yaml:"esp"`
}

// Validate validates the application configuration.
func (cfg *Config) Validate() error {
	if err := cfg.Log.Validate(); err != nil {
		return fmt.Errorf("validate log: %w", err)
	}

	if err := cfg.GRPC.Validate(); err != nil {
		return fmt.Errorf("validate GRPC: %w", err)
	}

	if err := cfg.HTTP.Validate(); err != nil {
		return fmt.Errorf("validate HTTP: %w", err)
	}

	if err := cfg.PIC.Validate(); err != nil {
		return fmt.Errorf("validate PIC: %w", err)
	}

	if err := cfg.ESP.Validate(); err != nil {
		return fmt.Errorf("validate ESP: %w", err)
	}

	return nil
}

// DefaultConfig is the default configuration for the application.
var DefaultConfig = Config{
	GRPC: GRPCConfig{
		Server: grpc.Config{
			Enable: true,
		},
		Cloud: cloud.Config{
			Address: "localhost:50051",
			Token:   "add token here",
		},
	},
	HTTP: http.Config{
		EnableSwagger: true,
	},
	Log: log.Config{
		Level:     "info",
		Format:    "text",
		AddSource: false,
	},
	PIC: picserial.Config{
		Serial: serial.Config{
			Port:        "/dev/ttyUSB0",
			BaudRate:    9600,
			DataBits:    8,
			Parity:      "none",
			StopBits:    1,
			ReadTimeout: 1 * time.Second,
		},
	},
	ESP: espserial.Config{
		Serial: espserial.SerialConfig{
			Port:        "/dev/ttyUSB1",
			BaudRate:    9600,
			DataBits:    8,
			Parity:      "none",
			StopBits:    1,
			ReadTimeout: 1 * time.Second,
		},
	},
}

func init() {
	// Ensure the default config is valid
	if err := DefaultConfig.Validate(); err != nil {
		panic(err)
	}
}
