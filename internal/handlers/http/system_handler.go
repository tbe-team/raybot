package http

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/system"
)

type systemHandler struct {
	systemService system.Service
}

func newSystemHandler(systemService system.Service) *systemHandler {
	return &systemHandler{systemService: systemService}
}

func (h systemHandler) RestartApplication(ctx context.Context, _ gen.RestartApplicationRequestObject) (gen.RestartApplicationResponseObject, error) {
	if err := h.systemService.RestartApplication(ctx); err != nil {
		return nil, fmt.Errorf("system service restart application: %w", err)
	}

	return gen.RestartApplication204Response{}, nil
}
