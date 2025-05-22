package http

import (
	"context"

	"github.com/tbe-team/raybot/internal/build"
	"github.com/tbe-team/raybot/internal/handlers/http/gen"
)

type versionHandler struct{}

func newVersionHandler() *versionHandler {
	return &versionHandler{}
}

func (h versionHandler) GetVersion(_ context.Context, _ gen.GetVersionRequestObject) (gen.GetVersionResponseObject, error) {
	i := build.GetBuildInfo()
	return gen.GetVersion200JSONResponse(gen.Version{
		Version:   i.Version,
		BuildDate: i.Date,
		GoVersion: i.GoVersion,
	}), nil
}
