package http

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	configsvc "github.com/tbe-team/raybot/internal/services/config"
	"github.com/tbe-team/raybot/pkg/log"
	"github.com/tbe-team/raybot/pkg/xerror"
)

type configHandler struct {
	configService configsvc.Service
}

func newConfigHandler(configService configsvc.Service) *configHandler {
	return &configHandler{
		configService: configService,
	}
}

func (h configHandler) GetLogConfig(ctx context.Context, _ gen.GetLogConfigRequestObject) (gen.GetLogConfigResponseObject, error) {
	cfg, err := h.configService.GetLogConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config service get log config: %w", err)
	}

	return gen.GetLogConfig200JSONResponse(h.convertLogConfigToResponse(cfg)), nil
}

func (h configHandler) UpdateLogConfig(ctx context.Context, request gen.UpdateLogConfigRequestObject) (gen.UpdateLogConfigResponseObject, error) {
	var level slog.Level
	switch request.Body.Level {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		return nil, xerror.ValidationFailed(nil, "invalid log level")
	}

	var format log.Format
	switch request.Body.Format {
	case "JSON":
		format = log.FormatJSON
	case "TEXT":
		format = log.FormatText
	default:
		return nil, xerror.ValidationFailed(nil, "invalid log format")
	}

	cfg, err := h.configService.UpdateLogConfig(ctx, config.Log{
		Level:     level,
		Format:    format,
		AddSource: request.Body.AddSource,
	})
	if err != nil {
		return nil, fmt.Errorf("config service update log config: %w", err)
	}

	return gen.UpdateLogConfig200JSONResponse(h.convertLogConfigToResponse(cfg)), nil
}

func (h configHandler) GetHardwareConfig(ctx context.Context, _ gen.GetHardwareConfigRequestObject) (gen.GetHardwareConfigResponseObject, error) {
	cfg, err := h.configService.GetHardwareConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config service get hardware config: %w", err)
	}

	return gen.GetHardwareConfig200JSONResponse(h.convertHardwareConfigToResponse(cfg)), nil
}

func (h configHandler) UpdateHardwareConfig(ctx context.Context, request gen.UpdateHardwareConfigRequestObject) (gen.UpdateHardwareConfigResponseObject, error) {

	//nolint:gosec
	espSerial := config.Serial{
		Port:        request.Body.Esp.Serial.Port,
		BaudRate:    request.Body.Esp.Serial.BaudRate,
		DataBits:    uint8(request.Body.Esp.Serial.DataBits),
		Parity:      request.Body.Esp.Serial.Parity,
		StopBits:    float32(request.Body.Esp.Serial.StopBits),
		ReadTimeout: time.Duration(request.Body.Esp.Serial.ReadTimeout * 1e9),
	}

	//nolint:gosec
	picSerial := config.Serial{
		Port:        request.Body.Pic.Serial.Port,
		BaudRate:    request.Body.Pic.Serial.BaudRate,
		DataBits:    uint8(request.Body.Pic.Serial.DataBits),
		Parity:      request.Body.Pic.Serial.Parity,
		StopBits:    float32(request.Body.Pic.Serial.StopBits),
		ReadTimeout: time.Duration(request.Body.Pic.Serial.ReadTimeout * 1e9),
	}

	cfg, err := h.configService.UpdateHardwareConfig(ctx, config.Hardware{
		ESP: config.ESP{
			Serial: espSerial,
		},
		PIC: config.PIC{
			Serial: picSerial,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("config service update hardware config: %w", err)
	}

	return gen.UpdateHardwareConfig200JSONResponse(h.convertHardwareConfigToResponse(cfg)), nil
}

func (h configHandler) GetCloudConfig(ctx context.Context, _ gen.GetCloudConfigRequestObject) (gen.GetCloudConfigResponseObject, error) {
	cfg, err := h.configService.GetCloudConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config service get cloud config: %w", err)
	}

	return gen.GetCloudConfig200JSONResponse(h.convertCloudConfigToResponse(cfg)), nil
}

func (h configHandler) UpdateCloudConfig(ctx context.Context, request gen.UpdateCloudConfigRequestObject) (gen.UpdateCloudConfigResponseObject, error) {
	cfg, err := h.configService.UpdateCloudConfig(ctx, config.Cloud{
		Address: request.Body.Address,
		Token:   request.Body.Token,
	})
	if err != nil {
		return nil, fmt.Errorf("config service update cloud config: %w", err)
	}

	return gen.UpdateCloudConfig200JSONResponse(h.convertCloudConfigToResponse(cfg)), nil
}

func (h configHandler) GetGRPCConfig(ctx context.Context, _ gen.GetGRPCConfigRequestObject) (gen.GetGRPCConfigResponseObject, error) {
	cfg, err := h.configService.GetGRPCConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config service get grpc config: %w", err)
	}

	return gen.GetGRPCConfig200JSONResponse(h.convertGRPCConfigToResponse(cfg)), nil
}

func (h configHandler) UpdateGRPCConfig(ctx context.Context, request gen.UpdateGRPCConfigRequestObject) (gen.UpdateGRPCConfigResponseObject, error) {
	//nolint:gosec
	cfg, err := h.configService.UpdateGRPCConfig(ctx, config.GRPC{
		Port:   uint32(request.Body.Port),
		Enable: request.Body.Enable,
	})
	if err != nil {
		return nil, fmt.Errorf("config service update grpc config: %w", err)
	}

	return gen.UpdateGRPCConfig200JSONResponse(h.convertGRPCConfigToResponse(cfg)), nil
}

func (h configHandler) GetHTTPConfig(ctx context.Context, _ gen.GetHTTPConfigRequestObject) (gen.GetHTTPConfigResponseObject, error) {
	cfg, err := h.configService.GetHTTPConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config service get http config: %w", err)
	}

	return gen.GetHTTPConfig200JSONResponse(h.convertHTTPConfigToResponse(cfg)), nil
}

func (h configHandler) UpdateHTTPConfig(ctx context.Context, request gen.UpdateHTTPConfigRequestObject) (gen.UpdateHTTPConfigResponseObject, error) {
	//nolint:gosec
	cfg, err := h.configService.UpdateHTTPConfig(ctx, config.HTTP{
		Port:    uint32(request.Body.Port),
		Swagger: request.Body.Swagger,
	})
	if err != nil {
		return nil, fmt.Errorf("config service update http config: %w", err)
	}

	return gen.UpdateHTTPConfig200JSONResponse(h.convertHTTPConfigToResponse(cfg)), nil
}

func (h configHandler) GetCargoConfig(ctx context.Context, _ gen.GetCargoConfigRequestObject) (gen.GetCargoConfigResponseObject, error) {
	cfg, err := h.configService.GetCargoConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config service get cargo config: %w", err)
	}

	return gen.GetCargoConfig200JSONResponse(h.convertCargoConfigToResponse(cfg)), nil
}

func (h configHandler) UpdateCargoConfig(ctx context.Context, request gen.UpdateCargoConfigRequestObject) (gen.UpdateCargoConfigResponseObject, error) {
	//nolint:gosec
	cfg, err := h.configService.UpdateCargoConfig(ctx, config.Cargo{
		LiftPosition:  uint16(request.Body.LiftPosition),
		LowerPosition: uint16(request.Body.LowerPosition),
	})
	if err != nil {
		return nil, fmt.Errorf("config service update cargo config: %w", err)
	}

	return gen.UpdateCargoConfig200JSONResponse(h.convertCargoConfigToResponse(cfg)), nil
}

func (h configHandler) GetWifiConfig(ctx context.Context, _ gen.GetWifiConfigRequestObject) (gen.GetWifiConfigResponseObject, error) {
	cfg, err := h.configService.GetWifiConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config service get wifi config: %w", err)
	}

	return gen.GetWifiConfig200JSONResponse(h.convertWifiConfigToResponse(cfg)), nil
}

func (h configHandler) UpdateWifiConfig(ctx context.Context, request gen.UpdateWifiConfigRequestObject) (gen.UpdateWifiConfigResponseObject, error) {
	cfg, err := h.configService.UpdateWifiConfig(ctx, config.Wifi{
		AP: config.APConfig{
			Enable:   request.Body.Ap.Enable,
			SSID:     request.Body.Ap.Ssid,
			Password: request.Body.Ap.Password,
			IP:       request.Body.Ap.Ip,
		},
		STA: config.STAConfig{
			Enable:   request.Body.Sta.Enable,
			SSID:     request.Body.Sta.Ssid,
			Password: request.Body.Sta.Password,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("config service update wifi config: %w", err)
	}

	return gen.UpdateWifiConfig200JSONResponse(h.convertWifiConfigToResponse(cfg)), nil
}

func (configHandler) convertLogConfigToResponse(cfg config.Log) gen.LogConfig {
	return gen.LogConfig{
		Level:     cfg.Level.String(),
		Format:    cfg.Format.String(),
		AddSource: cfg.AddSource,
	}
}

func (h configHandler) convertHardwareConfigToResponse(cfg config.Hardware) gen.HardwareConfig {
	return gen.HardwareConfig{
		Esp: gen.ESPConfig{
			Serial: h.convertSerialConfigToResponse(cfg.ESP.Serial),
		},
		Pic: gen.PICConfig{
			Serial: h.convertSerialConfigToResponse(cfg.PIC.Serial),
		},
	}
}

func (configHandler) convertSerialConfigToResponse(cfg config.Serial) gen.SerialConfig {
	return gen.SerialConfig{
		Port:        cfg.Port,
		BaudRate:    cfg.BaudRate,
		DataBits:    int(cfg.DataBits),
		Parity:      cfg.Parity,
		StopBits:    float64(cfg.StopBits),
		ReadTimeout: int(cfg.ReadTimeout.Seconds()),
	}
}

func (configHandler) convertCloudConfigToResponse(cfg config.Cloud) gen.CloudConfig {
	return gen.CloudConfig{
		Address: cfg.Address,
		Token:   cfg.Token,
	}
}

func (configHandler) convertGRPCConfigToResponse(cfg config.GRPC) gen.GRPCConfig {
	return gen.GRPCConfig{
		Port:   int(cfg.Port),
		Enable: cfg.Enable,
	}
}

func (configHandler) convertHTTPConfigToResponse(cfg config.HTTP) gen.HTTPConfig {
	return gen.HTTPConfig{
		Port:    int(cfg.Port),
		Swagger: cfg.Swagger,
	}
}

func (configHandler) convertCargoConfigToResponse(cfg config.Cargo) gen.CargoConfig {
	return gen.CargoConfig{
		LiftPosition:  int(cfg.LiftPosition),
		LowerPosition: int(cfg.LowerPosition),
	}
}

func (configHandler) convertWifiConfigToResponse(cfg config.Wifi) gen.WifiConfig {
	ap := gen.APConfig{
		Enable:   cfg.AP.Enable,
		Ssid:     cfg.AP.SSID,
		Password: cfg.AP.Password,
		Ip:       cfg.AP.IP,
	}
	sta := gen.STAConfig{
		Enable:   cfg.STA.Enable,
		Ssid:     cfg.STA.SSID,
		Password: cfg.STA.Password,
	}

	return gen.WifiConfig{
		Ap:  ap,
		Sta: sta,
	}
}
