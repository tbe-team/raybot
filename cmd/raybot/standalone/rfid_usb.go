package standalone

import (
	"fmt"
	"sync"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/handlers/rfidusb"
)

func startRFIDUSB(app *application.Application, interruptChan <-chan any, readyWg *sync.WaitGroup) error {
	service := rfidusb.New(
		app.Log,
		app.EventBus,
		app.LocationService,
	)

	cleanup, err := service.Run(app.Context)
	if err != nil {
		return fmt.Errorf("error running RFID USB service: %w", err)
	}

	app.Log.Info("rfid usb service started")

	readyWg.Done()
	<-interruptChan

	app.Log.Debug("rfid usb service is shutting down")

	if err := cleanup(app.Context); err != nil {
		return fmt.Errorf("error cleaning up RFID USB service: %w", err)
	}

	app.Log.Debug("rfid usb service stopped")

	return nil
}
