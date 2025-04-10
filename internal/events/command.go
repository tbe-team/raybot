package events

import "github.com/maniartech/signals"

type CommandCreatedEvent struct {
	CommandID int64
}

var CommandCreatedSignal = signals.New[CommandCreatedEvent]()
