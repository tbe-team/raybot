package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type PICSerialCommandRepository interface {
	GetPICSerialCommand(ctx context.Context, id string) (model.PICSerialCommand, error)
	CreatePICSerialCommand(ctx context.Context, command model.PICSerialCommand) error
	DeletePICSerialCommand(ctx context.Context, id string) error
}
