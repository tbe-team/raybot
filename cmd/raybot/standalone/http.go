package standalone

import (
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/handlers/http"
)

func startHTTPService(app *application.Application, interruptChan <-chan any) error {
	service := http.New(
		app.Cfg.HTTP,
		app.Log,
		app.ConfigService,
		app.SystemService,
		app.DashboardDataService,
		app.PeripheralService,
		app.CommandService,
		app.ApperrorcodeService,
		app.EmergencyService,
	)

	cleanup, err := service.Run()
	if err != nil {
		return fmt.Errorf("error running HTTP service: %w", err)
	}

	app.Log.Info("http service started", slog.Any("http_cfg", app.Cfg.HTTP))

	<-interruptChan

	app.Log.Debug("http service is shutting down")

	if err := cleanup(app.Context); err != nil {
		return fmt.Errorf("error cleaning up HTTP service: %w", err)
	}

	app.Log.Debug("http service stopped")

	return nil
}
