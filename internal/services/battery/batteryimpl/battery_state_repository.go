package batteryimpl

import (
	"context"
	"sync"
	"time"

	"github.com/tbe-team/raybot/internal/services/battery"
)

type batteryStateRepository struct {
	battery battery.BatteryState
	mu      sync.RWMutex
}

func NewBatteryStateRepository() battery.BatteryStateRepository {
	return &batteryStateRepository{
		battery: battery.BatteryState{
			CellVoltages: []uint16{0, 0, 0, 0}, // avoid nil
		},
	}
}

func (r *batteryStateRepository) GetBatteryState(_ context.Context) (battery.BatteryState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.battery, nil
}

func (r *batteryStateRepository) UpdateBatteryState(_ context.Context, params battery.UpdateBatteryStateParams) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.battery = battery.BatteryState{
		Current:      params.Current,
		Temp:         params.Temp,
		Voltage:      params.Voltage,
		CellVoltages: params.CellVoltages,
		Percent:      params.Percent,
		Fault:        params.Fault,
		Health:       params.Health,
		UpdatedAt:    time.Now(),
	}
	return nil
}
