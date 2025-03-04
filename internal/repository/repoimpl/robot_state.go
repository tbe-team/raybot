package repoimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
)

type robotStateRepository struct {
	mu         sync.RWMutex
	robotState model.RobotState
}

func newRobotStateRepository() repository.RobotStateRepository {
	return &robotStateRepository{
		robotState: model.RobotState{},
	}
}

func (r *robotStateRepository) GetRobotState(_ context.Context) (model.RobotState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.robotState, nil
}

func (r *robotStateRepository) UpdateRobotState(_ context.Context, state model.RobotState) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.robotState = state
	return nil
}
