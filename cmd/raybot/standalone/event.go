package standalone

import (
	"fmt"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/handlers/event"
)

func startEventService(app *application.Application, interruptChan <-chan any) error {
	service := event.New(app.Log, app.AppConnectionService)

	cleanup, err := service.Run(app.Context)
	if err != nil {
		return fmt.Errorf("error running event service: %w", err)
	}

	app.Log.Info("event service started")

	<-interruptChan

	app.Log.Debug("event service is shutting down")

	if err := cleanup(app.Context); err != nil {
		return fmt.Errorf("error cleaning up event service: %w", err)
	}

	app.Log.Debug("event service stopped")

	return nil
}
