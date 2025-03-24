package repository

import (
	"context"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/sort"
)

type UpdateCommandParams struct {
	ID             string
	Status         model.CommandStatus
	SetStatus      bool
	Error          *string
	SetError       bool
	CompletedAt    *time.Time
	SetCompletedAt bool
}

type CommandRepository interface {
	ListCommands(ctx context.Context, db db.SQLDB, params paging.Params, sorts []sort.Sort) (paging.List[model.Command], error)
	GetCommandByStatusInProgress(ctx context.Context, db db.SQLDB) (model.Command, error)
	GetCommandByID(ctx context.Context, db db.SQLDB, id string) (model.Command, error)
	CreateCommand(ctx context.Context, db db.SQLDB, command model.Command) error
	UpdateCommand(ctx context.Context, db db.SQLDB, params UpdateCommandParams) (model.Command, error)
}
