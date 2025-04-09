package events

import "github.com/maniartech/signals"

type OpenCargoDoorEvent struct {
}

type CloseCargoDoorEvent struct {
}

var OpenCargoDoorSignal = signals.New[OpenCargoDoorEvent]()
var CloseCargoDoorSignal = signals.New[CloseCargoDoorEvent]()
