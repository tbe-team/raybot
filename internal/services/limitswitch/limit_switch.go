package limitswitch

import "context"

type UpdateLimitSwitchStateParams struct {
	ID      LimitSwitchID
	Pressed bool
}

type Service interface {
	UpdateLimitSwitchState(ctx context.Context, params UpdateLimitSwitchStateParams) error
}

type Repository interface {
	GetLimitSwitchState(ctx context.Context, id LimitSwitchID) (LimitSwitch, error)
	UpdateLimitSwitchState(ctx context.Context, id LimitSwitchID, pressed bool) error
}
