package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
)

const dropBoxPosition = 100

type DropBoxExecutor struct {
	commandRepo repository.CommandRepository
	picService  service.PICService
	dbProvider  db.Provider
	log         *slog.Logger
}

func NewDropBoxExecutor(
	commandRepo repository.CommandRepository,
	picService service.PICService,
	dbProvider db.Provider,
	log *slog.Logger,
) *DropBoxExecutor {
	return &DropBoxExecutor{
		commandRepo: commandRepo,
		picService:  picService,
		dbProvider:  dbProvider,
		log:         log,
	}
}

func (e DropBoxExecutor) Execute(ctx context.Context, command model.Command) error {
	if command.Type != model.CommandTypeDropBox {
		return fmt.Errorf("command type is not drop box")
	}

	params := service.CreateSerialCommandParams{
		Data: model.PICSerialCommandBatteryLiftMotorData{
			Enable:         true,
			TargetPosition: dropBoxPosition,
		},
	}
	if err := e.picService.CreateSerialCommand(ctx, params); err != nil {
		return fmt.Errorf("create serial command: %w", err)
	}

	return nil
}
