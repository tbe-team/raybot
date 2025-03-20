package handler

import (
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/service"
)

var _ gen.StrictServerInterface = (*APIHandler)(nil)

type APIHandler struct {
	*systemHandler
	*healthHandler
	*robotStateHandler
	*commandHandler
}

func NewAPIHandler(service service.Service) *APIHandler {
	return &APIHandler{
		systemHandler:     &systemHandler{systemService: service.SystemService()},
		healthHandler:     &healthHandler{},
		robotStateHandler: &robotStateHandler{robotStateService: service.RobotStateService()},
		commandHandler:    &commandHandler{commandService: service.CommandService()},
	}
}
