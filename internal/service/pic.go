package service

import (
	"context"

	"github.com/tbe-team/raybot/internal/model"
)

type CreateSerialCommandParams struct {
	Data model.PICSerialCommandData `validate:"required"`
}

type ProcessSerialCommandACKParams struct {
	ID      string `validate:"required"`
	Success bool
}

type PICService interface {
	CreateSerialCommand(ctx context.Context, params CreateSerialCommandParams) error
	ProcessSerialCommandACK(ctx context.Context, params ProcessSerialCommandACKParams) error
}
