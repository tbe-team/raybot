package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/sort"
)

type CommandRepository interface {
	ListCommands(ctx context.Context, db db.SQLDB, params paging.Params, sorts []sort.Sort) ([]model.Command, error)
	GetCommandByStatusInProgress(ctx context.Context, db db.SQLDB) (model.Command, error)
	CreateCommand(ctx context.Context, db db.SQLDB, command model.Command) (model.Command, error)
}
