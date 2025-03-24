package repository

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type ESPSerialCommandRepository interface {
	GetESPSerialCommand(ctx context.Context, id string) (model.ESPSerialCommand, error)
	CreateESPSerialCommand(ctx context.Context, command model.ESPSerialCommand) error
	DeleteESPSerialCommand(ctx context.Context, id string) error
}
