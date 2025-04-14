package command

import (
	"context"
	"time"

	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/sort"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var (
	ErrCommandNotFound         = xerror.NotFound(nil, "command.notFound", "command not found")
	ErrNoNextExecutableCommand = xerror.NotFound(nil, "command.noNextExecutable", "no next executable command")
)

type CreateCommandParams struct {
	Source Source `validate:"enum"`
	Inputs Inputs `validate:"required"`
}

type GetCommandByIDParams struct {
	CommandID int64 `validate:"required,min=1"`
}

type ListCommandsParams struct {
	PagingParams paging.Params `validate:"required"`
	Sorts        []sort.Sort   `validate:"sort=type status source created_at updated_at completed_at"`
	Statuses     []Status      `validate:"dive,enum"`
}

type ExecuteCreatedCommandParams struct {
	CommandID int64
}

type Service interface {
	GetCommandByID(ctx context.Context, params GetCommandByIDParams) (Command, error)
	GetCurrentProcessingCommand(ctx context.Context) (Command, error)
	ListCommands(ctx context.Context, params ListCommandsParams) (paging.List[Command], error)
	CreateCommand(ctx context.Context, params CreateCommandParams) (Command, error)
	// CancelCommand(ctx context.Context, id int64) error
	// RetryCommand(ctx context.Context, id int64) error

	ExecuteCreatedCommand(ctx context.Context, params ExecuteCreatedCommandParams) error
}

type UpdateCommandParams struct {
	ID             int64
	Status         Status
	SetStatus      bool
	Error          *string
	SetError       bool
	StartedAt      *time.Time
	SetStartedAt   bool
	CompletedAt    *time.Time
	SetCompletedAt bool
	UpdatedAt      time.Time
}

type Repository interface {
	ListCommands(ctx context.Context, params ListCommandsParams) (paging.List[Command], error)
	GetNextExecutableCommand(ctx context.Context) (Command, error)
	GetCurrentProcessingCommand(ctx context.Context) (Command, error)
	CommandProcessingExists(ctx context.Context) (bool, error)
	GetCommandByID(ctx context.Context, id int64) (Command, error)
	CreateCommand(ctx context.Context, command Command) (Command, error)
	UpdateCommand(ctx context.Context, params UpdateCommandParams) (Command, error)
}
