package led

import (
	"context"
	"time"
)

type BlinkSystemLedParams struct {
	Duration time.Duration
}

type BlinkAlertLedParams struct {
	Duration time.Duration
}

type Service interface {
	SetSystemLedOn(ctx context.Context) error
	SetSystemLedOff(ctx context.Context) error
	BlinkSystemLed(ctx context.Context, params BlinkSystemLedParams) error

	SetAlertLedOn(ctx context.Context) error
	SetAlertLedOff(ctx context.Context) error
	BlinkAlertLed(ctx context.Context, params BlinkAlertLedParams) error
}

type LedsOutput struct {
	SystemLedState      State
	SystemLedConnection Connection
	AlertLedState       State
	AlertLedConnection  Connection
}

type Repository interface {
	GetLeds(ctx context.Context) (LedsOutput, error)
	UpdateSystemLedState(ctx context.Context, state State) error
	UpdateAlertLedState(ctx context.Context, state State) error
	UpdateSystemLedConnection(ctx context.Context, connection Connection) error
	UpdateAlertLedConnection(ctx context.Context, connection Connection) error
}
