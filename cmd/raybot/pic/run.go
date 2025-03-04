package pic

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/controller/picserial"
)

func Start(app *application.Application) error {
	picSerialService, err := picserial.NewPICSerialService(app.Cfg.PIC, app.Service)
	if err != nil {
		return fmt.Errorf("failed to create PIC serial service: %w", err)
	}

	cleanup, err := picSerialService.Run(app.Context())
	if err != nil {
		return fmt.Errorf("failed to run PIC serial service: %w", err)
	}

	app.Log.Debug("PIC serial service is running")

	app.CleanupManager.Add(func(ctx context.Context) error {
		app.Log.Debug("PIC serial service is shutting down")
		if err := cleanup(ctx); err != nil {
			return fmt.Errorf("failed to cleanup PIC serial service: %w", err)
		}

		app.Log.Debug("PIC serial service shut down complete")

		return nil
	})

	return nil
}
