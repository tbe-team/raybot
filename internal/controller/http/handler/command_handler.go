package handler

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/controller/http/converter"
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/sort"
)

type commandHandler struct {
	commandService service.CommandService
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
			return nil, fmt.Errorf("invalid sort: %w", err)
		}
	}

	params := service.ListCommandsParams{
		PagingParams: paging.Params{
			Page:     page,
			PageSize: pageSize,
		},
		Sorts: sorts,
	}

	commands, err := h.commandService.ListCommands(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("command service list commands: %w", err)
	}

	res := make([]gen.CommandResponse, len(commands.Items))
	for i, command := range commands.Items {
		response, err := converter.ConvertCommandToResponse(command)
		if err != nil {
			return nil, fmt.Errorf("convert command to response: %w", err)
		}
		res[i] = response
	}

	return gen.ListCommands200JSONResponse{
		TotalItems: int(commands.TotalItems),
		Items:      res,
	}, nil
}

func (h commandHandler) GetCurrentProcessingCommand(ctx context.Context, _ gen.GetCurrentProcessingCommandRequestObject) (gen.GetCurrentProcessingCommandResponseObject, error) {
	command, err := h.commandService.GetCurrentProcessingCommand(ctx)
	if err != nil {
		return nil, fmt.Errorf("command service get current processing command: %w", err)
	}

	response, err := converter.ConvertCommandToResponse(command)
	if err != nil {
		return nil, fmt.Errorf("convert command to response: %w", err)
	}

	return gen.GetCurrentProcessingCommand200JSONResponse(response), nil
}
