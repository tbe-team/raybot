package robotstate

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type Service struct {
	robotStateRepo repository.RobotStateRepository
	dbProvider     db.Provider
}

func NewService(
	robotStateRepo repository.RobotStateRepository,
	dbProvider db.Provider,
) *Service {
	return &Service{
		robotStateRepo: robotStateRepo,
		dbProvider:     dbProvider,
	}
}

func (s Service) GetRobotState(ctx context.Context) (model.RobotState, error) {
	return s.robotStateRepo.GetRobotState(ctx, s.dbProvider.DB())
}
