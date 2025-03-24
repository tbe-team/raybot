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

func (h *CommandHandler) MoveToLocation(ctx context.Context, req *commandv1.MoveToLocationRequest) (*commandv1.MoveToLocationResponse, error) {
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

func (h *CommandHandler) LiftBox(ctx context.Context, _ *commandv1.LiftBoxRequest) (*commandv1.LiftBoxResponse, error) {
	params := service.CreateCommandParams{
		Source:      model.CommandSourceCloud,
		CommandType: model.CommandTypeLiftBox,
		Inputs:      model.CommandLiftBoxInputs{},
	}
	command, err := h.commandService.CreateCommand(ctx, params)
	if err != nil {
		return nil, err
	}

	return &commandv1.LiftBoxResponse{CommandId: command.ID}, nil
}

func (h *CommandHandler) DropBox(ctx context.Context, _ *commandv1.DropBoxRequest) (*commandv1.DropBoxResponse, error) {
	params := service.CreateCommandParams{
		Source:      model.CommandSourceCloud,
		CommandType: model.CommandTypeDropBox,
		Inputs:      model.CommandDropBoxInputs{},
	}
	command, err := h.commandService.CreateCommand(ctx, params)
	if err != nil {
		return nil, err
	}

	return &commandv1.DropBoxResponse{CommandId: command.ID}, nil
}
