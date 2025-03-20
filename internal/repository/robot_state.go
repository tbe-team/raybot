package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type RobotStateRepository interface {
	GetRobotState(ctx context.Context, db db.SQLDB) (model.RobotState, error)
	UpdateRobotState(ctx context.Context, db db.SQLDB, state model.RobotState) error
}
