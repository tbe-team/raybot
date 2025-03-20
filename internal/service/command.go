package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/sort"
)

type ListCommandsParams struct {
	PagingParams paging.Params `validate:"required"`
	Sorts        []sort.Sort   `validate:"required"`
}

type CommandService interface {
	ListCommands(ctx context.Context, params ListCommandsParams) ([]model.Command, error)
	GetCurrentProcessingCommand(ctx context.Context) (model.Command, error)
}
