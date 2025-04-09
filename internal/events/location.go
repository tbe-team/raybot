package events

import "github.com/maniartech/signals"

type UpdateLocationEvent struct {
	CurrentLocation string
}

var UpdateLocationSignal = signals.New[UpdateLocationEvent]()
