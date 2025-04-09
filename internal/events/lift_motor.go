package events

import "github.com/maniartech/signals"

type UpdateLiftMotorStateEvent struct {
	TargetPosition uint16
	Enable         bool
}

var UpdateLiftMotorStateSignal = signals.New[UpdateLiftMotorStateEvent]()
