package application

import (
	"context"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/repository/repoimpl"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/service/serviceimpl"
	"github.com/tbe-team/raybot/pkg/log"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Application struct {
	CfgManager config.Manager

	Service service.Service

	Log *slog.Logger

	CleanupManager *CleanupManager

	ctx context.Context
}

func (a *Application) Context() context.Context {
	return a.ctx
}

type CleanupFunc func() error

func New(cfgManager config.Manager) (*Application, CleanupFunc, error) {
	// Set UTC timezone
	time.Local = time.UTC
	// Create context
	ctx := context.Background()

	logger := log.NewLogger(cfgManager.GetConfig().Log)
	slog.SetDefault(logger)

	// Setup repository
	repo := repoimpl.New()

	// Setup service
	validator := validator.New()
	service := serviceimpl.New(cfgManager, repo, validator)

	// Setup application
	app := &Application{
		CfgManager:     cfgManager,
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
