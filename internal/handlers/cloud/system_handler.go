package cloud

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/durationpb"

	sysv1 "github.com/tbe-team/raybot-api/sys/v1"
	"github.com/tbe-team/raybot/internal/build"
	"github.com/tbe-team/raybot/internal/services/system"
)

type systemHandler struct {
	sysv1.UnimplementedSysServiceServer
	systemService system.Service
}

func newSystemHandler(systemService system.Service) sysv1.SysServiceServer {
	return &systemHandler{
		systemService: systemService,
	}
}

func (h systemHandler) GetSysInfo(ctx context.Context, _ *sysv1.GetSysInfoRequest) (*sysv1.GetSysInfoResponse, error) {
	info, err := h.systemService.GetInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}

	return &sysv1.GetSysInfoResponse{
		LocalIp:     info.LocalIP,
		CpuUsage:    float32(info.CPUUsage),
		MemoryUsage: float32(info.MemoryUsage),
		TotalMemory: info.TotalMemory,
		Uptime:      durationpb.New(info.Uptime),
	}, nil
}

func (h systemHandler) GetVersion(_ context.Context, _ *sysv1.GetVersionRequest) (*sysv1.GetVersionResponse, error) {
	i := build.GetBuildInfo()
	return &sysv1.GetVersionResponse{
		Version:   i.Version,
		Date:      i.Version,
		GoVersion: i.GoVersion,
	}, nil
}

func (h systemHandler) Ping(_ context.Context, _ *sysv1.PingRequest) (*sysv1.PingResponse, error) {
	return &sysv1.PingResponse{}, nil
}

func (h systemHandler) Reboot(ctx context.Context, _ *sysv1.RebootRequest) (*sysv1.RebootResponse, error) {
	if err := h.systemService.Reboot(ctx); err != nil {
		return nil, fmt.Errorf("failed to reboot: %w", err)
	}
	return &sysv1.RebootResponse{}, nil
}

func (h systemHandler) StopEmergency(ctx context.Context, _ *sysv1.StopEmergencyRequest) (*sysv1.StopEmergencyResponse, error) {
	if err := h.systemService.StopEmergency(ctx); err != nil {
		return nil, fmt.Errorf("failed to stop emergency: %w", err)
	}
	return &sysv1.StopEmergencyResponse{}, nil
}
