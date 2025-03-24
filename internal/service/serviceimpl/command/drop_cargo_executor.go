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

const dropCargoPosition = 100

type DropCargoExecutor struct {
	commandRepo repository.CommandRepository
	picService  service.PICService
	dbProvider  db.Provider
	log         *slog.Logger
}

func NewDropCargoExecutor(
	commandRepo repository.CommandRepository,
	picService service.PICService,
	dbProvider db.Provider,
	log *slog.Logger,
) *DropCargoExecutor {
	return &DropCargoExecutor{
		commandRepo: commandRepo,
		picService:  picService,
		dbProvider:  dbProvider,
		log:         log,
	}
}

func (e DropCargoExecutor) Execute(ctx context.Context, command model.Command) error {
	if command.Type != model.CommandTypeDropCargo {
		return fmt.Errorf("command type is not drop cargo")
	}

	params := service.CreateSerialCommandParams{
		Data: model.PICSerialCommandBatteryLiftMotorData{
			Enable:         true,
			TargetPosition: dropCargoPosition,
		},
	}
	if err := e.picService.CreateSerialCommand(ctx, params); err != nil {
		return fmt.Errorf("create serial command: %w", err)
	}

	return nil
}
