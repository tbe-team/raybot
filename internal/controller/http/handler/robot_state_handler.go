package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/tbe-team/raybot/internal/controller/http/converter"
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/service"
)

type robotStateHandler struct {
	robotService service.RobotService
}

func (h robotStateHandler) GetRobotState(ctx context.Context, _ gen.GetRobotStateRequestObject) (gen.GetRobotStateResponseObject, error) {
	state, err := h.robotService.GetRobotState(ctx)
	log.Printf("state: %+v", state)
	if err != nil {
		return nil, fmt.Errorf("robot service get robot state: %w", err)
	}

	return gen.GetRobotState200JSONResponse(converter.ToRobotStateResponse(state)), nil
}
