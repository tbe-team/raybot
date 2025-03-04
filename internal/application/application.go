package application

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/tbe-team/raybot/internal/controller/picserial"
	"github.com/tbe-team/raybot/internal/repository/repoimpl"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/service/serviceimpl"
	"github.com/tbe-team/raybot/pkg/log"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Config struct {
	PIC picserial.Config `yaml:"pic"`
	Log log.Config       `yaml:"log"`
}

// RegisterFlags registers flags for the application.
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	cfg.PIC.RegisterFlags(f)
	cfg.Log.RegisterFlags(f)
}

// Validate validates the application configuration.
func (cfg *Config) Validate() error {
	if err := cfg.PIC.Validate(); err != nil {
		return fmt.Errorf("validate PIC: %w", err)
	}

	if err := cfg.Log.Validate(); err != nil {
		return fmt.Errorf("validate log: %w", err)
	}

	return nil
}

type Application struct {
	Cfg Config

	Service service.Service

	Log *slog.Logger

	CleanupManager *CleanupManager

	ctx context.Context
}

func (a *Application) Context() context.Context {
	return a.ctx
}

type CleanupFunc func() error

func New(cfg Config) (*Application, CleanupFunc, error) {
	// Set UTC timezone
	time.Local = time.UTC
	// Create context
	ctx := context.Background()

	logger := log.NewLogger(cfg.Log)
	slog.SetDefault(logger)

	// Setup repository
	repo := repoimpl.New()

	// Setup service
	validator := validator.New()
	service := serviceimpl.New(repo, validator)

	// Setup application
	app := &Application{
		Cfg:            cfg,
		Service:        service,
		Log:            logger,
		CleanupManager: NewCleanupManager(),
		ctx:            ctx,
	}

	// cleanup function
	cleanup := func() error {
		return app.CleanupManager.Cleanup(app.ctx)
	}

	return app, cleanup, nil
}

// LoadConfig loads configuration from file and command line flags
func LoadConfig(configPath string) (Config, error) {
	var cfg Config

	// Parse command line flags
	f := flag.NewFlagSet("raybot", flag.ContinueOnError)
	cfg.RegisterFlags(f)

	// Load from YAML file if provided
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return cfg, fmt.Errorf("read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return cfg, fmt.Errorf("parse config file: %w", err)
		}
	}

	if err := f.Parse(os.Args[1:]); err != nil {
		return cfg, fmt.Errorf("parse flags: %w", err)
	}

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		return cfg, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}
