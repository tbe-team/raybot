package handler

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/controller/http/converter"
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/service"
)

type robotStateHandler struct {
	robotStateService service.RobotStateService
}

func (h robotStateHandler) GetRobotState(ctx context.Context, _ gen.GetRobotStateRequestObject) (gen.GetRobotStateResponseObject, error) {
	state, err := h.robotStateService.GetRobotState(ctx)
	if err != nil {
		return nil, fmt.Errorf("robot state service get robot state: %w", err)
	}

	return gen.GetRobotState200JSONResponse(converter.ConvertRobotStateToResponse(state)), nil
}
