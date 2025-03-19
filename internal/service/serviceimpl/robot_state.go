package serviceimpl

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type RobotStateService struct {
	robotStateRepo repository.RobotStateRepository
	dbProvider     db.Provider
}

func NewRobotStateService(
	robotStateRepo repository.RobotStateRepository,
	dbProvider db.Provider,
) *RobotStateService {
	return &RobotStateService{
		robotStateRepo: robotStateRepo,
		dbProvider:     dbProvider,
	}
}

func (s RobotStateService) GetRobotState(ctx context.Context) (model.RobotState, error) {
	return s.robotStateRepo.GetRobotState(ctx, s.dbProvider.DB())
}
