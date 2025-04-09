package command

import (
	"context"

	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/sort"
)

type CreateCommandParams struct {
	Source Source `validate:"enum"`
	Inputs Inputs `validate:"required"`
}

type ListCommandsParams struct {
	PagingParams paging.Params `validate:"required"`
	Sorts        []sort.Sort   `validate:"sort=type status source created_at updated_at completed_at"`
	Statuses     []Status      `validate:"dive,enum"`
}

type Service interface {
	ListCommands(ctx context.Context, params ListCommandsParams) (paging.List[Command], error)
	CreateCommand(ctx context.Context, params CreateCommandParams) (Command, error)
	// CancelCommand(ctx context.Context, id int64) error
	// RetryCommand(ctx context.Context, id int64) error
}

type Repository interface {
	ListCommands(ctx context.Context, params ListCommandsParams) (paging.List[Command], error)
	CreateCommand(ctx context.Context, command Command) (Command, error)
}
