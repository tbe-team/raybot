package repoimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type CargoRepository struct {
	queries *sqlc.Queries
}

func NewCargoRepository(queries *sqlc.Queries) *CargoRepository {
	return &CargoRepository{queries: queries}
}

func (r CargoRepository) UpdateCargo(ctx context.Context, db db.SQLDB, params repository.UpdateCargoParams) (model.Cargo, error) {
	arg := sqlc.CargoUpdateParams{
		IsOpen:            boolToInt64(params.IsOpen),
		SetIsOpen:         params.SetIsOpen,
		QrCode:            params.QRCode,
		SetQrCode:         params.SetQRCode,
		BottomDistance:    int64(params.BottomDistance),
		SetBottomDistance: params.SetBottomDistance,
		UpdatedAt:         params.UpdatedAt.Format(time.RFC3339),
	}
	row, err := r.queries.CargoUpdate(ctx, db, arg)
	if err != nil {
		return model.Cargo{}, fmt.Errorf("queries update cargo: %w", err)
	}

	return r.convertRowToCargo(row)
}

func (r CargoRepository) UpdateCargoDoorMotor(ctx context.Context, db db.SQLDB, params repository.UpdateCargoDoorMotorParams) (model.CargoDoorMotor, error) {
	arg := sqlc.CargoDoorMotorUpdateParams{
		Direction:    int64(params.Direction),
		SetDirection: params.SetDirection,
		Speed:        int64(params.Speed),
		SetSpeed:     params.SetSpeed,
		IsRunning:    boolToInt64(params.IsRunning),
		SetIsRunning: params.SetIsRunning,
		Enabled:      boolToInt64(params.Enabled),
		SetEnabled:   params.SetEnabled,
		UpdatedAt:    params.UpdatedAt.Format(time.RFC3339),
	}
	row, err := r.queries.CargoDoorMotorUpdate(ctx, db, arg)
	if err != nil {
		return model.CargoDoorMotor{}, fmt.Errorf("queries update cargo door motor: %w", err)
	}

	return r.convertRowToCargoDoorMotor(row)
}

func (r CargoRepository) convertRowToCargo(row sqlc.Cargo) (model.Cargo, error) {
	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return model.Cargo{}, fmt.Errorf("parse updated at: %w", err)
	}

	//nolint:gosec
	return model.Cargo{
		IsOpen:         row.IsOpen == 1,
		QRCode:         row.QrCode,
		BottomDistance: uint16(row.BottomDistance),
		UpdatedAt:      updatedAt,
	}, nil
}

func (r CargoRepository) convertRowToCargoDoorMotor(row sqlc.CargoDoorMotor) (model.CargoDoorMotor, error) {
	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return model.CargoDoorMotor{}, fmt.Errorf("parse updated at: %w", err)
	}

	//nolint:gosec
	return model.CargoDoorMotor{
		Direction: model.CargoDoorMotorDirection(row.Direction),
		Speed:     uint8(row.Speed),
		IsRunning: row.IsRunning == 1,
		Enabled:   row.Enabled == 1,
		UpdatedAt: updatedAt,
	}, nil
}
