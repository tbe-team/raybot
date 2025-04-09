package events

import "github.com/maniartech/signals"

type MoveDirection uint8

const (
	MoveDirectionForward MoveDirection = iota
	MoveDirectionBackward
)

type UpdateDriveMotorStateEvent struct {
	Direction MoveDirection
	Speed     uint8
	Enable    bool
}

var UpdateDriveMotorStateSignal = signals.New[UpdateDriveMotorStateEvent]()
