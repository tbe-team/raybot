package handler

import (
	"context"

	commandv1 "github.com/tbe-team/raybot/internal/controller/grpc/gen/command/v1"
	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
)

type CommandHandler struct {
	commandv1.UnimplementedCommandServiceServer

	commandService service.CommandService
}

func NewCommandHandler(commandService service.CommandService) *CommandHandler {
	return &CommandHandler{
		commandService: commandService,
	}
}

func (h CommandHandler) MoveToLocation(ctx context.Context, req *commandv1.MoveToLocationRequest) (*commandv1.MoveToLocationResponse, error) {
	params := service.CreateCommandParams{
		Source:      model.CommandSourceCloud,
		CommandType: model.CommandTypeMoveToLocation,
		Inputs: model.CommandMoveToLocationInputs{
			Location: req.Location,
		},
	}
	command, err := h.commandService.CreateCommand(ctx, params)
	if err != nil {
		return nil, err
	}

	return &commandv1.MoveToLocationResponse{CommandId: command.ID}, nil
}

func (h CommandHandler) LiftCargo(ctx context.Context, _ *commandv1.LiftCargoRequest) (*commandv1.LiftCargoResponse, error) {
	params := service.CreateCommandParams{
		Source:      model.CommandSourceCloud,
		CommandType: model.CommandTypeLiftCargo,
		Inputs:      model.CommandLiftCargoInputs{},
	}
	command, err := h.commandService.CreateCommand(ctx, params)
	if err != nil {
		return nil, err
	}

	return &commandv1.LiftCargoResponse{CommandId: command.ID}, nil
}

func (h CommandHandler) DropCargo(ctx context.Context, _ *commandv1.DropCargoRequest) (*commandv1.DropCargoResponse, error) {
	params := service.CreateCommandParams{
		Source:      model.CommandSourceCloud,
		CommandType: model.CommandTypeDropCargo,
		Inputs:      model.CommandDropCargoInputs{},
	}
	command, err := h.commandService.CreateCommand(ctx, params)
	if err != nil {
		return nil, err
	}

	return &commandv1.DropCargoResponse{CommandId: command.ID}, nil
}

func (h CommandHandler) OpenCargo(ctx context.Context, _ *commandv1.OpenCargoRequest) (*commandv1.OpenCargoResponse, error) {
	params := service.CreateCommandParams{
		Source:      model.CommandSourceCloud,
		CommandType: model.CommandTypeOpenCargo,
		Inputs:      model.CommandOpenCargoInputs{},
	}
	command, err := h.commandService.CreateCommand(ctx, params)
	if err != nil {
		return nil, err
	}

	return &commandv1.OpenCargoResponse{CommandId: command.ID}, nil
}

func (h CommandHandler) CloseCargo(ctx context.Context, _ *commandv1.CloseCargoRequest) (*commandv1.CloseCargoResponse, error) {
	params := service.CreateCommandParams{
		Source:      model.CommandSourceCloud,
		CommandType: model.CommandTypeCloseCargo,
		Inputs:      model.CommandCloseCargoInputs{},
	}
	command, err := h.commandService.CreateCommand(ctx, params)
	if err != nil {
		return nil, err
	}

	return &commandv1.CloseCargoResponse{CommandId: command.ID}, nil
}
