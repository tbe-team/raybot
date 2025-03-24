package repository

import (
	"context"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type UpdateCargoParams struct {
	IsOpen            bool
	SetIsOpen         bool
	QRCode            string
	SetQRCode         string
	BottomDistance    uint16
	SetBottomDistance uint16
	UpdatedAt         time.Time
}

type UpdateCargoDoorMotorParams struct {
	Direction    model.CargoDoorMotorDirection
	SetDirection bool
	Speed        uint8
	SetSpeed     bool
	IsRunning    bool
	SetIsRunning bool
	Enabled      bool
	SetEnabled   bool
	UpdatedAt    time.Time
}

type CargoRepository interface {
	UpdateCargo(ctx context.Context, db db.SQLDB, params UpdateCargoParams) (model.Cargo, error)
	UpdateCargoDoorMotor(ctx context.Context, db db.SQLDB, params UpdateCargoDoorMotorParams) (model.CargoDoorMotor, error)
}
