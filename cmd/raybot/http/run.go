package http

import (
	"fmt"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/controller/http"
)

func Start(app *application.Application) error {
	httpService, err := http.NewHTTPService(app.CfgManager.GetConfig().HTTP, app.Service, app.Log)
	if err != nil {
		return fmt.Errorf("failed to create HTTP service: %w", err)
	}

	cleanup, err := httpService.Run()
	if err != nil {
		return fmt.Errorf("failed to run HTTP service: %w", err)
	}

	app.CleanupManager.Add(cleanup)

	return nil
}
