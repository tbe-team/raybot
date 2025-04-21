package http

import (
	"context"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
)

type healthHandler struct {
}

func newHealthHandler() *healthHandler {
	return &healthHandler{}
}

func (h healthHandler) GetHealth(_ context.Context, _ gen.GetHealthRequestObject) (gen.GetHealthResponseObject, error) {
	return gen.GetHealth200JSONResponse(gen.HealthResponse{
		Status: "ok",
	}), nil
}
