package cargoimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type cargoRepository struct {
	db      db.DB
	queries *sqlc.Queries
}

func NewCargoRepository(db db.DB, queries *sqlc.Queries) cargo.Repository {
	return &cargoRepository{
		db:      db,
		queries: queries,
	}
}

func (r cargoRepository) GetCargo(ctx context.Context) (cargo.Cargo, error) {
	row, err := r.queries.CargoGet(ctx, r.db)
	if err != nil {
		return cargo.Cargo{}, fmt.Errorf("failed to get cargo: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return cargo.Cargo{}, fmt.Errorf("failed to parse updated at: %w", err)
	}

	//nolint:gosec
	return cargo.Cargo{
		IsOpen:         row.IsOpen == 1,
		QRCode:         row.QrCode,
		BottomDistance: uint16(row.BottomDistance),
		UpdatedAt:      updatedAt,
	}, nil
}

func (r cargoRepository) GetCargoDoorMotorState(ctx context.Context) (cargo.DoorMotorState, error) {
	row, err := r.queries.CargoDoorMotorGet(ctx, r.db)
	if err != nil {
		return cargo.DoorMotorState{}, fmt.Errorf("failed to get cargo door motor state: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return cargo.DoorMotorState{}, fmt.Errorf("failed to parse updated at: %w", err)
	}

	direction := cargo.DirectionOpen
	if row.Direction == 1 {
		direction = cargo.DirectionClose
	}

	//nolint:gosec
	return cargo.DoorMotorState{
		Direction: direction,
		Speed:     uint8(row.Speed),
		IsRunning: row.IsRunning == 1,
		Enabled:   row.Enabled == 1,
		UpdatedAt: updatedAt,
	}, nil
}

func (r cargoRepository) UpdateCargoDoor(ctx context.Context, params cargo.UpdateCargoDoorParams) error {
	if _, err := r.queries.CargoUpdateIsOpen(ctx, r.db, sqlc.CargoUpdateIsOpenParams{
		IsOpen:    boolToInt64(params.IsOpen),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("failed to update cargo door: %w", err)
	}

	return nil
}

func (r cargoRepository) UpdateCargoQRCode(ctx context.Context, params cargo.UpdateCargoQRCodeParams) error {
	if _, err := r.queries.CargoUpdateQRCode(ctx, r.db, sqlc.CargoUpdateQRCodeParams{
		QrCode:    params.QRCode,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("failed to update cargo qr code: %w", err)
	}

	return nil
}

func (r cargoRepository) UpdateCargoBottomDistance(ctx context.Context, params cargo.UpdateCargoBottomDistanceParams) error {
	if _, err := r.queries.CargoUpdateBottomDistance(ctx, r.db, sqlc.CargoUpdateBottomDistanceParams{
		BottomDistance: int64(params.BottomDistance),
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("failed to update cargo bottom distance: %w", err)
	}

	return nil
}

func (r cargoRepository) UpdateCargoDoorMotorState(ctx context.Context, params cargo.UpdateCargoDoorMotorStateParams) error {
	var direction int64
	switch params.Direction {
	case cargo.DirectionOpen:
		direction = 0
	case cargo.DirectionClose:
		direction = 1
	default:
		return fmt.Errorf("invalid direction: %s", params.Direction)
	}

	if _, err := r.queries.CargoDoorMotorUpdate(ctx, r.db, sqlc.CargoDoorMotorUpdateParams{
		Direction: direction,
		Speed:     int64(params.Speed),
		IsRunning: boolToInt64(params.IsRunning),
		Enabled:   boolToInt64(params.Enabled),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("failed to update cargo door motor state: %w", err)
	}

	return nil
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
