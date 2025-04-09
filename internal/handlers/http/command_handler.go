package http

import (
	"context"
	"encoding/json"
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
	inputsJSON, err := req.Body.Inputs.MarshalJSON()
	if err != nil {
		return nil, xerror.ValidationFailed(err, "invalid inputs")
	}

	inputs, err := command.UnmarshalInputs(command.CommandType(req.Body.Type), inputsJSON)
	if err != nil {
		return nil, fmt.Errorf("unmarshal inputs: %w", err)
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

func (commandHandler) convertCommandToResponse(cmd command.Command) (gen.CommandResponse, error) {
	inputsJSON, err := json.Marshal(cmd.Inputs)
	if err != nil {
		return gen.CommandResponse{}, fmt.Errorf("marshal inputs: %w", err)
	}

	var inputs gen.CommandInputs
	err = json.Unmarshal(inputsJSON, &inputs)
	if err != nil {
		return gen.CommandResponse{}, fmt.Errorf("unmarshal inputs: %w", err)
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
