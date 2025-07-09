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

func (h systemHandler) GetSystemInfo(ctx context.Context, _ gen.GetSystemInfoRequestObject) (gen.GetSystemInfoResponseObject, error) {
	info, err := h.systemService.GetInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("system service get info: %w", err)
	}

	return gen.GetSystemInfo200JSONResponse(h.convertSysInfoToResponse(info)), nil
}

func (h systemHandler) GetSystemStatus(ctx context.Context, _ gen.GetSystemStatusRequestObject) (gen.GetSystemStatusResponseObject, error) {
	status, err := h.systemService.GetStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("system service get status: %w", err)
	}

	return gen.GetSystemStatus200JSONResponse(gen.SystemStatus{
		Status: status.String(),
	}), nil
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

func (systemHandler) convertSysInfoToResponse(info system.Info) gen.SystemInfo {
	return gen.SystemInfo{
		LocalIp:     info.LocalIP,
		CpuUsage:    float32(info.CPUUsage),
		MemoryUsage: float32(info.MemoryUsage),
		TotalMemory: float32(info.TotalMemory),
		Uptime:      float32(info.Uptime.Seconds()),
	}
}
