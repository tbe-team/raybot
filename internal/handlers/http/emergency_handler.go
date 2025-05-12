package http

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/emergency"
)

type emergencyHandler struct {
	emergencyService emergency.Service
}

func newEmergencyHandler(emergencyService emergency.Service) *emergencyHandler {
	return &emergencyHandler{
		emergencyService: emergencyService,
	}
}

func (h emergencyHandler) GetEmergencyState(ctx context.Context, _ gen.GetEmergencyStateRequestObject) (gen.GetEmergencyStateResponseObject, error) {
	state, err := h.emergencyService.GetEmergencyState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get emergency state: %w", err)
	}

	return gen.GetEmergencyState200JSONResponse(gen.EmergencyState{
		Locked: state.Locked,
	}), nil
}

func (h emergencyHandler) ResumeEmergency(ctx context.Context, _ gen.ResumeEmergencyRequestObject) (gen.ResumeEmergencyResponseObject, error) {
	err := h.emergencyService.ResumeOperation(ctx)
	if err != nil {
		return nil, fmt.Errorf("resume emergency: %w", err)
	}

	return gen.ResumeEmergency204Response{}, nil
}

func (h emergencyHandler) StopEmergency(ctx context.Context, _ gen.StopEmergencyRequestObject) (gen.StopEmergencyResponseObject, error) {
	err := h.emergencyService.StopOperation(ctx)
	if err != nil {
		return nil, fmt.Errorf("stop emergency: %w", err)
	}

	return gen.StopEmergency204Response{}, nil
}
