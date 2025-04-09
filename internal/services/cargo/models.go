package cargo

import (
	"fmt"
	"time"
)

type Cargo struct {
	IsOpen         bool
	QRCode         string
	BottomDistance uint16
	UpdatedAt      time.Time
}

type DoorMotorState struct {
	Direction DoorDirection
	Speed     uint8
	IsRunning bool
	Enabled   bool
	UpdatedAt time.Time
}

type DoorDirection string

func (d DoorDirection) String() string {
	return string(d)
}

func (d DoorDirection) Validate() error {
	switch d {
	case DirectionOpen, DirectionClose:
		return nil
	default:
		return fmt.Errorf("invalid direction: %s", d)
	}
}

const (
	DirectionOpen  DoorDirection = "OPEN"
	DirectionClose DoorDirection = "CLOSE"
)
