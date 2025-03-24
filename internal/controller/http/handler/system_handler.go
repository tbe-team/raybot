package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/tbe-team/raybot/internal/controller/http/converter"
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/service"
)

type systemHandler struct {
	systemService service.SystemService
}

func (h systemHandler) GetSystemConfig(ctx context.Context, _ gen.GetSystemConfigRequestObject) (gen.GetSystemConfigResponseObject, error) {
	cfg, err := h.systemService.GetSystemConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("system service get system config: %w", err)
	}

	return gen.GetSystemConfig200JSONResponse(converter.ConvertSystemConfigToResponse(cfg)), nil
}

func (h systemHandler) UpdateSystemConfig(ctx context.Context, request gen.UpdateSystemConfigRequestObject) (gen.UpdateSystemConfigResponseObject, error) {
	params := service.UpdateSystemConfigParams{
		LogConfig: service.LogConfig{
			Level:     request.Body.Log.Level,
			Format:    request.Body.Log.Format,
			AddSource: request.Body.Log.AddSource,
		},
		PICConfig: service.PICConfig{
			Serial: service.SerialConfig{
				Port:        request.Body.Pic.Serial.Port,
				BaudRate:    request.Body.Pic.Serial.BaudRate,
				DataBits:    request.Body.Pic.Serial.DataBits,
				StopBits:    request.Body.Pic.Serial.StopBits,
				Parity:      request.Body.Pic.Serial.Parity,
				ReadTimeout: time.Duration(request.Body.Pic.Serial.ReadTimeout) * time.Second,
			},
		},
		ESPConfig: service.ESPConfig{
			Serial: service.SerialConfig{
				Port:        request.Body.Esp.Serial.Port,
				BaudRate:    request.Body.Esp.Serial.BaudRate,
				DataBits:    request.Body.Esp.Serial.DataBits,
				StopBits:    request.Body.Esp.Serial.StopBits,
				Parity:      request.Body.Esp.Serial.Parity,
				ReadTimeout: time.Duration(request.Body.Esp.Serial.ReadTimeout) * time.Second,
			},
		},
		GRPCConfig: service.GRPCConfig{
			Server: service.GRPCServerConfig{
				Enable: request.Body.Grpc.Server.Enable,
			},
			Cloud: service.CloudConfig{
				Address: request.Body.Grpc.Cloud.Address,
			},
		},
		HTTPConfig: service.HTTPConfig{
			EnableSwagger: request.Body.Http.EnableSwagger,
		},
	}
	cfg, err := h.systemService.UpdateSystemConfig(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("system service update system config: %w", err)
	}

	return gen.UpdateSystemConfig200JSONResponse(converter.ConvertSystemConfigToResponse(cfg)), nil
}

func (h systemHandler) RestartApplication(ctx context.Context, _ gen.RestartApplicationRequestObject) (gen.RestartApplicationResponseObject, error) {
	if err := h.systemService.RestartApplication(ctx); err != nil {
		return nil, fmt.Errorf("system service restart application: %w", err)
	}

	return gen.RestartApplication204Response{}, nil
}
