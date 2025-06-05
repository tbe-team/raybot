package http

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/limitswitch"
)

type stateHandler struct {
	limitSwitchService limitswitch.Service
}

func newStateHandler(limitSwitchService limitswitch.Service) *stateHandler {
	return &stateHandler{
		limitSwitchService: limitSwitchService,
	}
}

func (h stateHandler) GetLimitSwitchState(ctx context.Context, _ gen.GetLimitSwitchStateRequestObject) (gen.GetLimitSwitchStateResponseObject, error) {
	state, err := h.limitSwitchService.GetLimitSwitchState(ctx)
	if err != nil {
		return nil, fmt.Errorf("limit switch service get limit switch state: %w", err)
	}

	return gen.GetLimitSwitchState200JSONResponse{
		LimitSwitch1: gen.LimitSwitch{
			Pressed:   state.LimitSwitch1.Pressed,
			UpdatedAt: state.LimitSwitch1.UpdatedAt,
		},
	}, nil
}
