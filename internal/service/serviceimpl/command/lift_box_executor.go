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

const liftBoxPosition = 0

type LiftBoxExecutor struct {
	commandRepo repository.CommandRepository
	picService  service.PICService
	dbProvider  db.Provider
	log         *slog.Logger
}

func NewLiftBoxExecutor(
	commandRepo repository.CommandRepository,
	picService service.PICService,
	dbProvider db.Provider,
	log *slog.Logger,
) *LiftBoxExecutor {
	return &LiftBoxExecutor{
		commandRepo: commandRepo,
		picService:  picService,
		dbProvider:  dbProvider,
		log:         log,
	}
}

func (e LiftBoxExecutor) Execute(ctx context.Context, command model.Command) error {
	if command.Type != model.CommandTypeLiftBox {
		return fmt.Errorf("command type is not lift box")
	}

	//nolint:gosec
	params := service.CreateSerialCommandParams{
		Data: model.PICSerialCommandBatteryLiftMotorData{
			Enable:         true,
			TargetPosition: liftBoxPosition,
		},
	}
	if err := e.picService.CreateSerialCommand(ctx, params); err != nil {
		return fmt.Errorf("create serial command: %w", err)
	}

	return nil
}
