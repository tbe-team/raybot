package distancesensor

import "time"

//nolint:revive
type DistanceSensorState struct {
	FrontDistance uint16
	BackDistance  uint16
	DownDistance  uint16
	UpdatedAt     time.Time
}
