package handler

import (
	"context"
	"fmt"

	commandv1 "github.com/tbe-team/raybot/internal/controller/grpc/gen/command/v1"
	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/pkg/xerror"
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

func (h CommandHandler) CreateCommand(ctx context.Context, req *commandv1.CreateCommandRequest) (*commandv1.CreateCommandResponse, error) {
	commandType, err := h.convertReqCommandTypeToCommandType(req.Type)
	if err != nil {
		return nil, err
	}

	inputs, err := h.convertReqPayloadToCommandInputs(req.Payload)
	if err != nil {
		return nil, err
	}

	params := service.CreateCommandParams{
		Source:      model.CommandSourceCloud,
		CommandType: commandType,
		Inputs:      inputs,
	}
	cmd, err := h.commandService.CreateCommand(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("command service create command: %w", err)
	}

	return &commandv1.CreateCommandResponse{
		Id: cmd.ID,
	}, nil
}

func (CommandHandler) convertReqCommandTypeToCommandType(reqCommandType commandv1.CreateCommandRequest_Type) (model.CommandType, error) {
	switch reqCommandType {
	case commandv1.CreateCommandRequest_TYPE_MOVE_TO_LOCATION:
		return model.CommandTypeMoveToLocation, nil
	case commandv1.CreateCommandRequest_TYPE_LIFT_BOX:
		return model.CommandTypeLiftBox, nil
	case commandv1.CreateCommandRequest_TYPE_DROP_BOX:
		return model.CommandTypeDropBox, nil
	default:
		return 0, xerror.ValidationFailed(nil, fmt.Sprintf("invalid command type: %s", reqCommandType))
	}
}

func (CommandHandler) convertReqPayloadToCommandInputs(payload any) (model.CommandInputs, error) {
	switch payload := payload.(type) {
	case *commandv1.CreateCommandRequest_MoveToLocation:
		return model.CommandMoveToLocationInputs{
			Location: payload.MoveToLocation.Location,
		}, nil
	case *commandv1.CreateCommandRequest_LiftBox:
		return model.CommandLiftBoxInputs{}, nil
	case *commandv1.CreateCommandRequest_DropBox:
		return model.CommandDropBoxInputs{}, nil
	default:
		return nil, xerror.ValidationFailed(nil, fmt.Sprintf("invalid payload: %v", payload))
	}
}
