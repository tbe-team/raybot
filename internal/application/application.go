package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/controller/espserial"
	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
	"github.com/tbe-team/raybot/internal/pubsub"
	"github.com/tbe-team/raybot/internal/repository/repoimpl"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/service/serviceimpl"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/file"
	"github.com/tbe-team/raybot/pkg/log"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Application struct {
	CfgManager config.Manager

	PICSerialClient serial.Client
	ESPSerialClient espserial.Client

	Service service.Service
	PubSub  pubsub.PubSub

	Log *slog.Logger

	CleanupManager *CleanupManager

	ctx context.Context
}

func (a *Application) Context() context.Context {
	return a.ctx
}

type CleanupFunc func() error

func New() (*Application, CleanupFunc, error) {
	// Create context
	ctx := context.Background()

	path, err := NewPath()
	if err != nil {
		return nil, nil, fmt.Errorf("create path: %w", err)
	}

	fileClient, err := file.NewLocalFileClient(".")
	if err != nil {
		return nil, nil, fmt.Errorf("create file client: %w", err)
	}

	cfgManager, err := config.NewManager(fileClient, path.ConfigPath(), slog.Default())
	if err != nil {
		return nil, nil, fmt.Errorf("create config manager: %w", err)
	}

	logger := log.NewLogger(cfgManager.GetConfig().Log)
	slog.SetDefault(logger)

	// Setup repository
	dbProvider, err := db.NewProvider(db.Config{
		DBPath: path.DBPath(),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create db provider: %w", err)
	}

	// Auto migrate the database
	if err := dbProvider.AutoMigrate(); err != nil {
		return nil, nil, fmt.Errorf("failed to auto migrate the database: %w", err)
	}

	// Setup pubSub
	pubSub := pubsub.New(logger)

	// Setup repository
	repo := repoimpl.New()

	// Setup serial client
	picSerialClient, err := serial.NewClient(cfgManager.GetConfig().PIC.Serial, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create serial client: %w", err)
	}
	espSerialClient, err := espserial.NewClient(cfgManager.GetConfig().ESP.Serial, logger)
	if err != nil {
		logger.Error("failed to create esp serial client", slog.Any("error", err))
		// return nil, nil, fmt.Errorf("failed to create esp serial client: %w", err)
	}

	// Setup service
	validator := validator.New()
	service := serviceimpl.New(
		cfgManager,
		picSerialClient,
		espSerialClient,
		repo,
		pubSub,
		dbProvider,
		validator,
		logger,
	)

	// Setup application
	app := &Application{
		CfgManager:      cfgManager,
		PICSerialClient: picSerialClient,
		ESPSerialClient: espSerialClient,
		PubSub:          pubSub,
		Service:         service,
		Log:             logger,
		CleanupManager:  NewCleanupManager(),
		ctx:             ctx,
	}

	// cleanup function
	cleanup := func() error {
		if err := app.CleanupManager.Cleanup(app.ctx); err != nil {
			return fmt.Errorf("cleanup manager cleanup failed: %w", err)
		}

		if err := app.PICSerialClient.Stop(); err != nil {
			return fmt.Errorf("failed to close pic serial client: %w", err)
		}

		if app.ESPSerialClient != nil {
			if err := app.ESPSerialClient.Stop(); err != nil {
				return fmt.Errorf("failed to close esp serial client: %w", err)
			}
		}

		if err := dbProvider.Close(); err != nil {
			return fmt.Errorf("failed to close db provider: %w", err)
		}

		return nil
	}

	return app, cleanup, nil
}
