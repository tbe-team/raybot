package rfid

import (
	"fmt"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/controller/rfid"
)

func Start(app *application.Application) error {
	rfidService, err := rfid.NewRFIDService(app.Service, app.Log)
	if err != nil {
		return fmt.Errorf("failed to create rfid service: %w", err)
	}

	cleanup, err := rfidService.Run(app.Context())
	if err != nil {
		return fmt.Errorf("failed to run rfid service: %w", err)
	}

	app.CleanupManager.Add(cleanup)

	return nil
}
