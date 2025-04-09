package distancesensor

import "context"

type UpdateDistanceSensorStateParams struct {
	FrontDistance uint16
	BackDistance  uint16
	DownDistance  uint16
}

type Service interface {
	UpdateDistanceSensorState(ctx context.Context, params UpdateDistanceSensorStateParams) error
}

//nolint:revive
type DistanceSensorStateRepository interface {
	GetDistanceSensorState(ctx context.Context) (DistanceSensorState, error)
	UpdateDistanceSensorState(ctx context.Context, params UpdateDistanceSensorStateParams) error
}
