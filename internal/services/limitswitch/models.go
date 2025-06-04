package limitswitch

import "time"

type LimitSwitch struct {
	Pressed   bool
	UpdatedAt time.Time
}

//nolint:revive
type LimitSwitchID uint8

const (
	LimitSwitchID1 LimitSwitchID = 1
)
