package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type SyncCargoDoorStateParams struct {
	IsOpen bool
}

type SyncCargoQRCodeParams struct {
	QRCode string
}

type SyncCargoBottomDistanceParams struct {
	BottomDistance uint16
}

type SyncCargoDoorMotorStateParams struct {
	Direction model.CargoDoorMotorDirection `validate:"enum"`
	Speed     uint8                         `validate:"min=0,max=100"`
	IsRunning bool
	Enabled   bool
}

type ProcessESPSerialCommandACKParams struct {
	ID      string `validate:"required"`
	Success bool
}

type CargoControlService interface {
	SyncCargoDoorState(ctx context.Context, params SyncCargoDoorStateParams) error
	SyncCargoQRCode(ctx context.Context, params SyncCargoQRCodeParams) error
	SyncCargoBottomDistance(ctx context.Context, params SyncCargoBottomDistanceParams) error
	SyncCargoDoorMotorState(ctx context.Context, params SyncCargoDoorMotorStateParams) error

	ProcessESPSerialCommandACK(ctx context.Context, params ProcessESPSerialCommandACKParams) error

	OpenCargoDoor(ctx context.Context) error
	CloseCargoDoor(ctx context.Context) error
}
