package standalone

import (
	"log"
	"sync"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/pkg/cmdutil"
)

func Run(configFilePath, dbPath string) {
	app, cleanup, err := application.New(configFilePath, dbPath)
	if err != nil {
		log.Fatalf("error creating application: %v", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			log.Fatalf("error cleaning up: %v", err)
		}
	}()

	interruptChan := cmdutil.InterruptChan()
	var wg sync.WaitGroup

	// We need to start the event service first to ensure that the event handlers are registered
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startEventService(app, interruptChan); err != nil {
			log.Fatalf("error starting event service: %v", err)
		}
	}()

	var hardwareWgReady sync.WaitGroup
	hardwareWgReady.Add(3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startPICSerial(app, interruptChan, &hardwareWgReady); err != nil {
			log.Fatalf("error starting PIC serial service: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startESPSerial(app, interruptChan, &hardwareWgReady); err != nil {
			log.Fatalf("error starting ESP serial service: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startRFIDUSB(app, interruptChan, &hardwareWgReady); err != nil {
			log.Fatalf("error starting RFID USB service: %v", err)
		}
	}()

	// Wait for all hardware components to ensure they are ready
	hardwareWgReady.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startJobs(app, interruptChan); err != nil {
			log.Fatalf("error starting job service: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startCloud(app, interruptChan); err != nil {
			log.Fatalf("error starting cloud service: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startHTTPService(app, interruptChan); err != nil {
			log.Fatalf("error starting HTTP service: %v", err)
		}
	}()

	wg.Wait()
}
