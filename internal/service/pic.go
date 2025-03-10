package service

import "context"

type ProcessSerialCommandACK struct {
	ID      string `validate:"required"`
	Success bool
}

type PICService interface {
	ProcessSerialCommandACK(ctx context.Context, params ProcessSerialCommandACK) error
}
