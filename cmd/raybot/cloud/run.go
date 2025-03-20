package cloud

import (
	"fmt"

	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/internal/controller/cloud"
)

func Start(app *application.Application) error {
	cloudService, err := cloud.NewCloudService(app.CfgManager.GetConfig().GRPC.Cloud, app.Service, app.Log)
	if err != nil {
		return fmt.Errorf("failed to create cloud service: %w", err)
	}

	cleanup, err := cloudService.Run()
	if err != nil {
		return fmt.Errorf("failed to run cloud service: %w", err)
	}

	app.CleanupManager.Add(cleanup)

	return nil
}
