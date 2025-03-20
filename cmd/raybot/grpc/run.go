package grpc

import (
	"fmt"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/controller/grpc"
)

func Start(app *application.Application) error {
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
