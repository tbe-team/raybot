package model

import (
	"fmt"
	"time"
)

type DriveMotorDirection uint8

func (s DriveMotorDirection) Validate() error {
	switch s {
	case DriveMotorDirectionForward, DriveMotorDirectionBackward:
		return nil
	default:
		return fmt.Errorf("invalid drive motor direction: %d", s)
	}
}

func (s DriveMotorDirection) String() string {
	return []string{"FORWARD", "BACKWARD"}[s]
}

const (
	DriveMotorDirectionForward DriveMotorDirection = iota
	DriveMotorDirectionBackward
)

type DriveMotor struct {
	Direction DriveMotorDirection
	Speed     uint8
	IsRunning bool
	Enabled   bool
	UpdatedAt time.Time
}
