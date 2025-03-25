package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type RobotStateService interface {
	GetRobotState(ctx context.Context) (model.RobotState, error)
}
