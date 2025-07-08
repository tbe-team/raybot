package led

import (
	"fmt"
	"log/slog"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type Mode int

const (
	ModeOff Mode = iota
	ModeOn
	ModeBlink
)

type State struct {
	Mode     Mode
	Interval time.Duration // Blink interval if Mode == ModeBlink
}

type Led struct {
	name    string
	pin     gpio.PinIO
	stateCh chan State
	stopCh  chan struct{}
	log     *slog.Logger
}

func New(name, pin string, log *slog.Logger) (*Led, error) {
	p := gpioreg.ByName(pin)
	if p == nil {
		return nil, fmt.Errorf("pin not found: %s", pin)
	}

	return &Led{
		name:    name,
		pin:     p,
		stateCh: make(chan State, 1),
		stopCh:  make(chan struct{}),
		log:     log,
	}, nil
}

func (l *Led) Set(state State) {
	l.stateCh <- state
}

func (l *Led) Stop() {
	close(l.stopCh)
}

func (l *Led) Run() {
	curr := State{Mode: ModeOff}
	t := time.NewTicker(time.Hour)
	defer t.Stop()

	isOn := false

	for {
		switch curr.Mode {
		case ModeOn:
			if !isOn {
				l.on()
				isOn = true
			}

		case ModeOff:
			if isOn {
				l.off()
				isOn = false
			}

		case ModeBlink:
			t = time.NewTicker(curr.Interval)
		}

		select {
		case newState := <-l.stateCh:
			curr = newState
			if curr.Mode == ModeBlink {
				t.Reset(curr.Interval)
			}

		case <-t.C:
			if curr.Mode == ModeBlink {
				if isOn {
					l.off()
					isOn = false
				} else {
					l.on()
					isOn = true
				}
			}

		case <-l.stopCh:
			l.off()
			return
		}
	}
}

func (l *Led) on() {
	if err := l.pin.Out(gpio.High); err != nil {
		l.log.Error("set pin to high",
			slog.String("name", l.name),
			slog.String("pin", l.pin.String()),
			slog.Any("error", err))
	}
}

func (l *Led) off() {
	if err := l.pin.Out(gpio.Low); err != nil {
		l.log.Error("set pin to low",
			slog.String("name", l.name),
			slog.String("pin", l.pin.String()),
			slog.Any("error", err))
	}
}
