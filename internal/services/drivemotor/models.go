package drivemotor

import (
	"fmt"
	"time"
)

type Direction string

func (d Direction) Validate() error {
	switch d {
	case DirectionForward, DirectionBackward:
		return nil
	default:
		return fmt.Errorf("invalid drive motor direction: %s", d)
	}
}

func (d Direction) String() string {
	return string(d)
}

const (
	DirectionForward  Direction = "FORWARD"
	DirectionBackward Direction = "BACKWARD"
)

//nolint:revive
type DriveMotorState struct {
	Direction Direction
	Speed     uint8
	IsRunning bool
	Enabled   bool
	UpdatedAt time.Time
}
