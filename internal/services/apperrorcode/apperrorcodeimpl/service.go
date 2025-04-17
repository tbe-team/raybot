package apperrorcodeimpl

import (
	"context"

	"github.com/tbe-team/raybot/internal/services/apperrorcode"
)

type service struct {
}

func NewService() apperrorcode.Service {
	return &service{}
}

func (s service) ListErrorCodes(_ context.Context) ([]apperrorcode.ErrorCode, error) {
	return GetAll(), nil
}
