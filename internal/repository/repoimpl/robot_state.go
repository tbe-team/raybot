package repoimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/model"
)

type RobotStateRepository struct {
	mu         sync.RWMutex
	robotState model.RobotState
}

func NewRobotStateRepository() *RobotStateRepository {
	return &RobotStateRepository{
		robotState: model.RobotState{},
	}
}

func (r *RobotStateRepository) GetRobotState(_ context.Context) (model.RobotState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.robotState, nil
}

func (r *RobotStateRepository) UpdateRobotState(_ context.Context, state model.RobotState) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.robotState = state
	return nil
}
