package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/sort"
)

type CreateCommandParams struct {
	Source      model.CommandSource `validate:"enum"`
	CommandType model.CommandType   `validate:"enum"`
	Inputs      model.CommandInputs `validate:"required"`
}

type ListCommandsParams struct {
	PagingParams paging.Params `validate:"required"`
	Sorts        []sort.Sort   `validate:"sort=type status source created_at completed_at"`
}

type ExecuteCommandParams struct {
	CommandID string `validate:"required"`
}

type CommandService interface {
	ListCommands(ctx context.Context, params ListCommandsParams) (paging.List[model.Command], error)
	GetCurrentProcessingCommand(ctx context.Context) (model.Command, error)
	CreateCommand(ctx context.Context, params CreateCommandParams) (model.Command, error)
	ExecuteCommand(ctx context.Context, params ExecuteCommandParams) error
}
