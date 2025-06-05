package limitswitch

import "context"

type UpdateLimitSwitchByIDParams struct {
	ID      LimitSwitchID
	Pressed bool
}

type GetLimitSwitchStateOutput struct {
	LimitSwitch1 LimitSwitch
}

type Service interface {
	GetLimitSwitchState(ctx context.Context) (GetLimitSwitchStateOutput, error)
	UpdateLimitSwitchByID(ctx context.Context, params UpdateLimitSwitchByIDParams) error
}

type Repository interface {
	GetLimitSwitchByID(ctx context.Context, id LimitSwitchID) (LimitSwitch, error)
	UpdateLimitSwitchByID(ctx context.Context, id LimitSwitchID, pressed bool) error
}
