package http

import (
	"context"
	"fmt"
	"strings"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/sort"
	"github.com/tbe-team/raybot/pkg/xerror"
)

type commandHandler struct {
	commandService command.Service
}

func newCommandHandler(commandService command.Service) *commandHandler {
	return &commandHandler{
		commandService: commandService,
	}
}

//nolint:revive
func (h commandHandler) GetCommandById(ctx context.Context, request gen.GetCommandByIdRequestObject) (gen.GetCommandByIdResponseObject, error) {
	cmd, err := h.commandService.GetCommandByID(ctx, command.GetCommandByIDParams{
		CommandID: int64(request.CommandId),
	})
	if err != nil {
		return nil, fmt.Errorf("get command by id: %w", err)
	}

	res, err := h.convertCommandToResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("convert command to response: %w", err)
	}

	return gen.GetCommandById200JSONResponse(res), nil
}

func (h commandHandler) GetCurrentProcessingCommand(ctx context.Context, _ gen.GetCurrentProcessingCommandRequestObject) (gen.GetCurrentProcessingCommandResponseObject, error) {
	cmd, err := h.commandService.GetCurrentProcessingCommand(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current processing command: %w", err)
	}

	res, err := h.convertCommandToResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("convert command to response: %w", err)
	}

	return gen.GetCurrentProcessingCommand200JSONResponse(res), nil
}

func (h commandHandler) ListCommands(ctx context.Context, req gen.ListCommandsRequestObject) (gen.ListCommandsResponseObject, error) {
	page := uint(1)
	pageSize := uint(10)
	if req.Params.Page != nil {
		page = *req.Params.Page
	}
	if req.Params.PageSize != nil {
		pageSize = *req.Params.PageSize
	}

	sorts := []sort.Sort{}
	var err error
	if req.Params.Sorts != nil {
		sorts, err = sort.NewListFromString(*req.Params.Sorts)
		if err != nil {
			return nil, xerror.ValidationFailed(err, "invalid sort")
		}
	}

	statuses := []command.Status{}
	if req.Params.Statuses != nil && len(*req.Params.Statuses) > 0 {
		stripped := strings.TrimSpace(*req.Params.Statuses)
		if stripped == "" {
			return nil, xerror.ValidationFailed(nil, "invalid statuses")
		}
		ss := strings.Split(stripped, ",")
		for _, s := range ss {
			statuses = append(statuses, command.Status(s))
		}
	}

	commands, err := h.commandService.ListCommands(ctx, command.ListCommandsParams{
		PagingParams: paging.NewParams(paging.Page(page), paging.PageSize(pageSize)),
		Sorts:        sorts,
		Statuses:     statuses,
	})
	if err != nil {
		return nil, fmt.Errorf("list commands: %w", err)
	}

	res := make([]gen.CommandResponse, len(commands.Items))
	for i, cmd := range commands.Items {
		r, err := h.convertCommandToResponse(cmd)
		if err != nil {
			return nil, fmt.Errorf("convert command to response: %w", err)
		}
		res[i] = r
	}

	return gen.ListCommands200JSONResponse{
		TotalItems: int(commands.TotalItems),
		Items:      res,
	}, nil
}

func (h commandHandler) CreateCommand(ctx context.Context, req gen.CreateCommandRequestObject) (gen.CreateCommandResponseObject, error) {
	inputs, err := h.convertReqInputsToCommandInputs(req.Body.Type, req.Body.Inputs)
	if err != nil {
		return nil, xerror.ValidationFailed(err, "invalid inputs")
	}

	cmd, err := h.commandService.CreateCommand(ctx, command.CreateCommandParams{
		Source: command.SourceApp,
		Inputs: inputs,
	})
	if err != nil {
		return nil, fmt.Errorf("create command: %w", err)
	}

	res, err := h.convertCommandToResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("convert command to response: %w", err)
	}

	return gen.CreateCommand201JSONResponse(res), nil
}

func (h commandHandler) convertCommandToResponse(cmd command.Command) (gen.CommandResponse, error) {
	inputs, err := h.convertInputsToResponse(cmd.Inputs)
	if err != nil {
		return gen.CommandResponse{}, fmt.Errorf("convert inputs to response: %w", err)
	}

	return gen.CommandResponse{
		Id:          int(cmd.ID),
		Type:        cmd.Type.String(),
		Status:      cmd.Status.String(),
		Source:      cmd.Source.String(),
		Inputs:      inputs,
		Error:       cmd.Error,
		CompletedAt: cmd.CompletedAt,
		CreatedAt:   cmd.CreatedAt,
		UpdatedAt:   cmd.UpdatedAt,
	}, nil
}

func (commandHandler) convertInputsToResponse(inputs command.Inputs) (gen.CommandInputs, error) {
	var res gen.CommandInputs
	switch v := inputs.(type) {
	case *command.StopInputs:
		if err := res.FromStopInputs(gen.StopInputs{}); err != nil {
			return gen.CommandInputs{}, fmt.Errorf("from stop inputs: %w", err)
		}

	case *command.MoveToInputs:
		if err := res.FromMoveToInputs(gen.MoveToInputs{
			Location: v.Location,
		}); err != nil {
			return gen.CommandInputs{}, fmt.Errorf("from move to inputs: %w", err)
		}
	case *command.MoveForwardInputs:
		if err := res.FromMoveForwardInputs(gen.MoveForwardInputs{}); err != nil {
			return gen.CommandInputs{}, fmt.Errorf("from move forward inputs: %w", err)
		}
	case *command.MoveBackwardInputs:
		if err := res.FromMoveBackwardInputs(gen.MoveBackwardInputs{}); err != nil {
			return gen.CommandInputs{}, fmt.Errorf("from move backward inputs: %w", err)
		}
	case *command.CargoOpenInputs:
		if err := res.FromCargoOpenInputs(gen.CargoOpenInputs{}); err != nil {
			return gen.CommandInputs{}, fmt.Errorf("from cargo open inputs: %w", err)
		}
	case *command.CargoCloseInputs:
		if err := res.FromCargoCloseInputs(gen.CargoCloseInputs{}); err != nil {
			return gen.CommandInputs{}, fmt.Errorf("from cargo close inputs: %w", err)
		}
	case *command.CargoLiftInputs:
		if err := res.FromCargoLiftInputs(gen.CargoLiftInputs{}); err != nil {
			return gen.CommandInputs{}, fmt.Errorf("from cargo lift inputs: %w", err)
		}
	case *command.CargoLowerInputs:
		if err := res.FromCargoLowerInputs(gen.CargoLowerInputs{}); err != nil {
			return gen.CommandInputs{}, fmt.Errorf("from cargo lower inputs: %w", err)
		}
	case *command.CargoCheckQRInputs:
		if err := res.FromCargoCheckQRInputs(gen.CargoCheckQRInputs{
			QrCode: v.QRCode,
		}); err != nil {
			return gen.CommandInputs{}, fmt.Errorf("from cargo check qr inputs: %w", err)
		}
	default:
		return gen.CommandInputs{}, fmt.Errorf("unknown inputs type: %T", v)
	}

	return res, nil
}

func (commandHandler) convertReqInputsToCommandInputs(cmdType gen.CommandType, inputs gen.CommandInputs) (command.Inputs, error) {
	switch command.CommandType(cmdType) {
	case command.CommandTypeStop:
		return &command.StopInputs{}, nil

	case command.CommandTypeMoveTo:
		i, err := inputs.AsMoveToInputs()
		if err != nil {
			return nil, fmt.Errorf("as move to inputs: %w", err)
		}
		return &command.MoveToInputs{
			Location: i.Location,
		}, nil

	case command.CommandTypeMoveForward:
		return &command.MoveForwardInputs{}, nil

	case command.CommandTypeMoveBackward:
		return &command.MoveBackwardInputs{}, nil

	case command.CommandTypeCargoOpen:
		return &command.CargoOpenInputs{}, nil

	case command.CommandTypeCargoClose:
		return &command.CargoCloseInputs{}, nil

	case command.CommandTypeCargoLift:
		return &command.CargoLiftInputs{}, nil

	case command.CommandTypeCargoLower:
		return &command.CargoLowerInputs{}, nil

	case command.CommandTypeCargoCheckQR:
		return &command.CargoCheckQRInputs{}, xerror.NotImplemented(nil, "cargo.checkQR", "this command type is not implemented")
		// i, err := inputs.AsCargoCheckQRInputs()
		// if err != nil {
		// 	return nil, fmt.Errorf("as cargo check qr inputs: %w", err)
		// }
		// return &command.CargoCheckQRInputs{
		// 	QRCode: i.QrCode,
		// }, nil

	default:
		return nil, fmt.Errorf("unknown command type: %s", cmdType)
	}
}
