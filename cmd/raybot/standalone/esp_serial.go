package standalone

import (
	"fmt"
	"sync"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/handlers/espserial"
)

func startESPSerial(app *application.Application, interruptChan <-chan any, readyWg *sync.WaitGroup) error {
	service := espserial.New(
		app.Cfg.Hardware.ESP,
		app.Log,
		app.EventBus,
		app.ESPSerialClient,
		app.CargoService,
	)

	cleanup, err := service.Run(app.Context)
	if err != nil {
		return fmt.Errorf("error running ESP serial service: %w", err)
	}

	app.Log.Info("esp serial service started")

	readyWg.Done()
	<-interruptChan

	app.Log.Debug("esp serial service is shutting down")

	if err := cleanup(app.Context); err != nil {
		return fmt.Errorf("error cleaning up ESP serial service: %w", err)
	}

	app.Log.Debug("esp serial service stopped")

	return nil
}
