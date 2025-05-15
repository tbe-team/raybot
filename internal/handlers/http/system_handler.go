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

func (h systemHandler) RebootSystem(ctx context.Context, _ gen.RebootSystemRequestObject) (gen.RebootSystemResponseObject, error) {
	if err := h.systemService.Reboot(ctx); err != nil {
		return nil, fmt.Errorf("system service reboot: %w", err)
	}

	return gen.RebootSystem204Response{}, nil
}

func (h systemHandler) StopEmergency(ctx context.Context, _ gen.StopEmergencyRequestObject) (gen.StopEmergencyResponseObject, error) {
	if err := h.systemService.StopEmergency(ctx); err != nil {
		return nil, fmt.Errorf("system service stop emergency: %w", err)
	}

	return gen.StopEmergency204Response{}, nil
}
