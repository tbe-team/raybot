package standalone

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/controller/cloud"
	"github.com/tbe-team/raybot/internal/controller/espserial"
	"github.com/tbe-team/raybot/internal/controller/event"
	"github.com/tbe-team/raybot/internal/controller/grpc"
	"github.com/tbe-team/raybot/internal/controller/http"
	"github.com/tbe-team/raybot/internal/controller/picserial"
	"github.com/tbe-team/raybot/internal/controller/rfid"
	"github.com/tbe-team/raybot/pkg/cmdutil"
)

// Run starts all services in standalone mode
func Run() {
	app, cleanup, err := application.New()
	if err != nil {
		log.Printf("failed to create application: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := cleanup(); err != nil {
			log.Printf("failed to cleanup application: %v\n", err)
			os.Exit(1)
		}
	}()

	interruptChan := cmdutil.InterruptChan()

	// Ensure PIC serial service, ESP serial service, and RFID service are started before any other services
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runPIC(app); err != nil {
			log.Printf("failed to run PIC serial service: %v\n", err)
			os.Exit(1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runESP(app); err != nil {
			log.Printf("failed to run ESP serial service: %v\n", err)
			os.Exit(1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runRFID(app); err != nil {
			log.Printf("failed to run RFID service: %v\n", err)
			os.Exit(1)
		}
	}()

	wg.Wait()

	go func() {
		if err := runCloud(app); err != nil {
			log.Printf("failed to run cloud service: %v\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		if err := runEvent(app); err != nil {
			log.Printf("failed to run event service: %v\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		if err := runGRPC(app); err != nil {
			log.Printf("failed to run GRPC service: %v\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		if err := runHTTP(app); err != nil {
			log.Printf("failed to run HTTP service: %v\n", err)
			os.Exit(1)
		}
	}()

	<-interruptChan
}

func runPIC(app *application.Application) error {
	picSerialService, err := picserial.NewPICSerialService(app.CfgManager.GetConfig().PIC, app.PICSerialClient, app.Service, app.Log)
	if err != nil {
		return fmt.Errorf("failed to create PIC serial service: %w", err)
	}

	cleanup, err := picSerialService.Run(app.Context())
	if err != nil {
		return fmt.Errorf("failed to run PIC serial service: %w", err)
	}

	app.CleanupManager.Add(cleanup)

	return nil
}

func runESP(app *application.Application) error {
	espSerialService, err := espserial.New(app.CfgManager.GetConfig().ESP, app.ESPSerialClient, app.Service, app.Log)
	if err != nil {
		return fmt.Errorf("failed to create ESP serial service: %w", err)
	}

	cleanup, err := espSerialService.Run(app.Context())
	if err != nil {
		return fmt.Errorf("failed to run ESP serial service: %w", err)
	}

	app.CleanupManager.Add(cleanup)

	return nil
}

func runRFID(app *application.Application) error {
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

func runCloud(app *application.Application) error {
	cloudService, err := cloud.NewService(app.CfgManager.GetConfig().GRPC.Cloud, app.Service, app.Log)
	if err != nil {
		return fmt.Errorf("failed to create cloud service: %w", err)
	}

	cleanup, err := cloudService.Run(app.Context())
	if err != nil {
		return fmt.Errorf("failed to run cloud service: %w", err)
	}

	app.CleanupManager.Add(cleanup)

	return nil
}

func runEvent(app *application.Application) error {
	eventService := event.New(app.Service, app.PubSub, app.Log)

	cleanup, err := eventService.Run(app.Context())
	if err != nil {
		return fmt.Errorf("failed to run event service: %w", err)
	}

	app.CleanupManager.Add(cleanup)

	return nil
}

func runGRPC(app *application.Application) error {
	grpcService, err := grpc.NewGRPCService(app.CfgManager.GetConfig().GRPC.Server, app.Service, app.Log)
	if err != nil {
		return fmt.Errorf("failed to create GRPC service: %w", err)
	}

	cleanup, err := grpcService.Run()
	if err != nil {
		return fmt.Errorf("failed to run GRPC service: %w", err)
	}

	app.CleanupManager.Add(cleanup)

	return nil
}

func runHTTP(app *application.Application) error {
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
