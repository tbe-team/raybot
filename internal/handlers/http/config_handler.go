package http

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	configsvc "github.com/tbe-team/raybot/internal/services/config"
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
	cfg, err := h.convertLogConfigReqToModel(*request.Body)
	if err != nil {
		return nil, fmt.Errorf("convert log config request to model: %w", err)
	}

	cfg, err = h.configService.UpdateLogConfig(ctx, cfg)
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
		BottomDistanceHysteresis: config.CargoBottomDistanceHysteresis{
			LowerThreshold: uint16(request.Body.BottomDistanceHysteresis.LowerThreshold),
			UpperThreshold: uint16(request.Body.BottomDistanceHysteresis.UpperThreshold),
		},
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
		File: gen.LogFileHandler{
			Enable:        cfg.File.Enable,
			Path:          cfg.File.Path,
			RotationCount: cfg.File.RotationCount,
			Level:         cfg.File.Level.String(),
			Format:        cfg.File.Format.String(),
		},
		Console: gen.LogConsoleHandler{
			Enable: cfg.Console.Enable,
			Level:  cfg.Console.Level.String(),
			Format: cfg.Console.Format.String(),
		},
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
		BottomDistanceHysteresis: gen.CargoBottomDistanceHysteresis{
			LowerThreshold: int(cfg.BottomDistanceHysteresis.LowerThreshold),
			UpperThreshold: int(cfg.BottomDistanceHysteresis.UpperThreshold),
		},
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

func (h configHandler) convertLogConfigReqToModel(req gen.LogConfig) (config.Log, error) {
	fileLogLevel, err := h.convertLogLevelReqToModel(req.File.Level)
	if err != nil {
		return config.Log{}, err
	}

	fileLogFormat, err := h.convertLogFormatReqToModel(req.File.Format)
	if err != nil {
		return config.Log{}, err
	}

	consoleLogLevel, err := h.convertLogLevelReqToModel(req.Console.Level)
	if err != nil {
		return config.Log{}, err
	}

	consoleLogFormat, err := h.convertLogFormatReqToModel(req.Console.Format)
	if err != nil {
		return config.Log{}, err
	}

	return config.Log{
		File: config.LogFileHandler{
			Enable:        req.File.Enable,
			Path:          req.File.Path,
			RotationCount: req.File.RotationCount,
			Level:         fileLogLevel,
			Format:        fileLogFormat,
		},
		Console: config.LogConsoleHandler{
			Enable: req.Console.Enable,
			Level:  consoleLogLevel,
			Format: consoleLogFormat,
		},
	}, nil
}

func (configHandler) convertLogLevelReqToModel(level string) (slog.Level, error) {
	switch level {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, xerror.ValidationFailed(nil, "invalid log level")
	}
}

func (configHandler) convertLogFormatReqToModel(format string) (config.LogFormat, error) {
	switch format {
	case "JSON":
		return config.LogFormatJSON, nil
	case "TEXT":
		return config.LogFormatText, nil
	default:
		return config.LogFormatText, xerror.ValidationFailed(nil, "invalid log format")
	}
}
