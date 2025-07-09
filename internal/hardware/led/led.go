package led

import (
	"context"
	"errors"
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type Led struct {
	pin gpio.PinIO
}

func New(pin string) (*Led, error) {
	p := gpioreg.ByName(pin)
	if p == nil {
		return nil, fmt.Errorf("pin not found: %s", pin)
	}

	return &Led{
		pin: p,
	}, nil
}

func (l *Led) On() error {
	return l.pin.Out(gpio.High)
}

func (l *Led) Off() error {
	return l.pin.Out(gpio.Low)
}

// Blink blinks the led on and off at the given interval.
// If the context is done, the led will be turned off.
func (l *Led) Blink(ctx context.Context, interval time.Duration) error {
	isOn := l.pin.Read()

	toggleFunc := func() error {
		if isOn {
			return l.Off()
		}
		return l.On()
	}

	for {
		select {
		case <-ctx.Done():
			if err := l.Off(); err != nil {
				return errors.Join(err, ctx.Err())
			}
			return ctx.Err()

		case <-time.After(interval):
			if err := toggleFunc(); err != nil {
				return err
			}
			isOn = !isOn
		}
	}
}

func (l *Led) Stop() error {
	return l.pin.Halt()
}
