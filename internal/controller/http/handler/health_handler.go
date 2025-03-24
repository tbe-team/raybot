package handler

import (
	"context"

	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
)

type healthHandler struct{}

func (h healthHandler) GetHealth(_ context.Context, _ gen.GetHealthRequestObject) (gen.GetHealthResponseObject, error) {
	return gen.GetHealth200JSONResponse(gen.HealthResponse{
		Status: "ok",
	}), nil
}
