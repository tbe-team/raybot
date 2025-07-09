package led

import "time"

type Mode string

func (m Mode) String() string {
	return string(m)
}

const (
	ModeOff   Mode = "OFF"
	ModeOn    Mode = "ON"
	ModeBlink Mode = "BLINK"
)

type State struct {
	Mode      Mode
	UpdatedAt time.Time
}

type Connection struct {
	Connected       bool
	LastConnectedAt *time.Time
	Error           *string
}
