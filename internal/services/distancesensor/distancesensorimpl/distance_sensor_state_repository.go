package distancesensorimpl

import (
	"context"
	"sync"
	"time"

	"github.com/tbe-team/raybot/internal/services/distancesensor"
)

type distanceSensorStateRepository struct {
	distanceSensor distancesensor.DistanceSensorState
	mu             sync.RWMutex
}

func NewDistanceSensorStateRepository() distancesensor.DistanceSensorStateRepository {
	return &distanceSensorStateRepository{
		distanceSensor: distancesensor.DistanceSensorState{},
	}
}

func (r *distanceSensorStateRepository) GetDistanceSensorState(_ context.Context) (distancesensor.DistanceSensorState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.distanceSensor, nil
}

func (r *distanceSensorStateRepository) UpdateDistanceSensorState(_ context.Context, params distancesensor.UpdateDistanceSensorStateParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.distanceSensor = distancesensor.DistanceSensorState{
		FrontDistance: params.FrontDistance,
		BackDistance:  params.BackDistance,
		DownDistance:  params.DownDistance,
		UpdatedAt:     time.Now(),
	}

	return nil
}
