package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type LiftMotorRepository interface {
	GetLiftMotor(ctx context.Context, db db.SQLDB) (model.LiftMotor, error)
	UpdateLiftMotor(ctx context.Context, db db.SQLDB, liftMotor model.LiftMotor) error
}
