package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type RobotStateRepository interface {
	GetRobotState(ctx context.Context) (model.RobotState, error)
	UpdateRobotState(ctx context.Context, state model.RobotState) error
}
