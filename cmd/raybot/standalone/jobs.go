package standalone

import (
	"fmt"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/handlers/jobs"
)

func startJobs(app *application.Application, interruptChan <-chan any) error {
	service := jobs.New(app.Cfg.Cron, app.Log, app.CommandService)

	cleanup, err := service.Run()
	if err != nil {
		return fmt.Errorf("failed to run job service: %w", err)
	}

	app.Log.Info("job service started")

	<-interruptChan

	app.Log.Debug("job service is shutting down")

	if err := cleanup(app.Context); err != nil {
		return fmt.Errorf("failed to cleanup job service: %w", err)
	}

	app.Log.Debug("job service stopped")

	return nil
}
