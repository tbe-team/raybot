package standalone

import (
	"fmt"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/handlers/cloud"
)

func startCloud(app *application.Application, interruptChan <-chan any) error {
	if app.Cfg.Wifi.AP.Enable {
		return nil
	}

	service, err := cloud.New(app.Cfg.Cloud, app.Log, app.EventBus)
	if err != nil {
		return fmt.Errorf("failed to create cloud service: %w", err)
	}

	cleanup, err := service.Run(app.Context)
	if err != nil {
		return fmt.Errorf("failed to run cloud service: %w", err)
	}

	app.Log.Info("cloud service started")

	<-interruptChan

	app.Log.Debug("cloud service is shutting down")

	if err := cleanup(app.Context); err != nil {
		return fmt.Errorf("failed to cleanup cloud service: %w", err)
	}

	app.Log.Debug("cloud service stopped")

	return nil
}
