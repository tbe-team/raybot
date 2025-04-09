package liftmotor

import "time"

//nolint:revive
type LiftMotorState struct {
	CurrentPosition uint16
	TargetPosition  uint16
	IsRunning       bool
	Enabled         bool
	UpdatedAt       time.Time
}
