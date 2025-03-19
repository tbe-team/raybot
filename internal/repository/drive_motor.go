package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type DriveMotorRepository interface {
	GetDriveMotor(ctx context.Context, db db.SQLDB) (model.DriveMotor, error)
	UpdateDriveMotor(ctx context.Context, db db.SQLDB, driveMotor model.DriveMotor) error
}
