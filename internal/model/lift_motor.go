package model

import "time"

type LiftMotor struct {
	CurrentPosition uint16
	TargetPosition  uint16
	IsRunning       bool
	Enabled         bool
	UpdatedAt       time.Time
}
