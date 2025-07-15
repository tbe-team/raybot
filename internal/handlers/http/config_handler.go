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
		ReadTimeout: time.Duration(request.Body.Esp.Serial.ReadTimeout) * time.Second,
	}

	//nolint:gosec
	picSerial := config.Serial{
		Port:        request.Body.Pic.Serial.Port,
		BaudRate:    request.Body.Pic.Serial.BaudRate,
		DataBits:    uint8(request.Body.Pic.Serial.DataBits),
		Parity:      request.Body.Pic.Serial.Parity,
		StopBits:    float32(request.Body.Pic.Serial.StopBits),
		ReadTimeout: time.Duration(request.Body.Pic.Serial.ReadTimeout) * time.Second,
	}

	cfg, err := h.configService.UpdateHardwareConfig(ctx, config.Hardware{
		ESP: config.ESP{
			Serial:            espSerial,
			EnableACK:         request.Body.Esp.EnableAck,
			CommandACKTimeout: time.Duration(request.Body.Esp.CommandAckTimeout) * time.Millisecond,
		},
		PIC: config.PIC{
			Serial:            picSerial,
			EnableACK:         request.Body.Pic.EnableAck,
			CommandACKTimeout: time.Duration(request.Body.Pic.CommandAckTimeout) * time.Millisecond,
		},
		Leds: config.Leds{
			System: config.Led{
				Pin: request.Body.Leds.System.Pin,
			},
			Alert: config.Led{
				Pin: request.Body.Leds.Alert.Pin,
			},
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
		Enable:  request.Body.Enable,
		Address: request.Body.Address,
		Token:   request.Body.Token,
	})
	if err != nil {
		return nil, fmt.Errorf("config service update cloud config: %w", err)
	}

	return gen.UpdateCloudConfig200JSONResponse(h.convertCloudConfigToResponse(cfg)), nil
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
			IP:       request.Body.Sta.Ip,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("config service update wifi config: %w", err)
	}

	return gen.UpdateWifiConfig200JSONResponse(h.convertWifiConfigToResponse(cfg)), nil
}

func (h configHandler) GetCommandConfig(ctx context.Context, _ gen.GetCommandConfigRequestObject) (gen.GetCommandConfigResponseObject, error) {
	cfg, err := h.configService.GetCommandConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config service get command config: %w", err)
	}

	return gen.GetCommandConfig200JSONResponse(h.convertCommandConfigToResponse(cfg)), nil
}

func (h configHandler) UpdateCommandConfig(ctx context.Context, req gen.UpdateCommandConfigRequestObject) (gen.UpdateCommandConfigResponseObject, error) {
	cfg, err := h.configService.UpdateCommandConfig(ctx, config.Command{
		CargoLift: config.CargoLift{
			StableReadCount: req.Body.CargoLift.StableReadCount,
		},
		CargoLower: config.CargoLower{
			StableReadCount: req.Body.CargoLower.StableReadCount,
			BottomObstacleTracking: config.ObstacleTracking{
				EnterDistance: req.Body.CargoLower.BottomObstacleTracking.EnterDistance,
				ExitDistance:  req.Body.CargoLower.BottomObstacleTracking.ExitDistance,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("config service update command config: %w", err)
	}

	return gen.UpdateCommandConfig200JSONResponse(h.convertCommandConfigToResponse(cfg)), nil
}

func (h configHandler) GetBatteryMonitoringConfig(ctx context.Context, _ gen.GetBatteryMonitoringConfigRequestObject) (gen.GetBatteryMonitoringConfigResponseObject, error) {
	cfg, err := h.configService.GetBatteryMonitoringConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config service get battery monitoring config: %w", err)
	}

	return gen.GetBatteryMonitoringConfig200JSONResponse(h.convertBatteryMonitoringConfigToResponse(cfg)), nil
}

func (h configHandler) UpdateBatteryMonitoringConfig(ctx context.Context, request gen.UpdateBatteryMonitoringConfigRequestObject) (gen.UpdateBatteryMonitoringConfigResponseObject, error) {
	cfg, err := h.configService.UpdateBatteryMonitoringConfig(ctx, config.BatteryMonitoring{
		BatteryVoltageLow: config.BatteryVoltageLow{
			Enable:    request.Body.VoltageLow.Enable,
			Threshold: request.Body.VoltageLow.Threshold,
		},
		BatteryVoltageHigh: config.BatteryVoltageHigh{
			Enable:    request.Body.VoltageHigh.Enable,
			Threshold: request.Body.VoltageHigh.Threshold,
		},
		BatteryCellVoltageHigh: config.BatteryCellVoltageHigh{
			Enable:    request.Body.CellVoltageHigh.Enable,
			Threshold: request.Body.CellVoltageHigh.Threshold,
		},
		BatteryCellVoltageLow: config.BatteryCellVoltageLow{
			Enable:    request.Body.CellVoltageLow.Enable,
			Threshold: request.Body.CellVoltageLow.Threshold,
		},
		BatteryCellVoltageDiff: config.BatteryCellVoltageDiff{
			Enable:    request.Body.CellVoltageDiff.Enable,
			Threshold: request.Body.CellVoltageDiff.Threshold,
		},
		BatteryCurrentHigh: config.BatteryCurrentHigh{
			Enable:    request.Body.CurrentHigh.Enable,
			Threshold: request.Body.CurrentHigh.Threshold,
		},
		BatteryTempHigh: config.BatteryTempHigh{
			Enable:    request.Body.TempHigh.Enable,
			Threshold: request.Body.TempHigh.Threshold,
		},
		BatteryPercentLow: config.BatteryPercentLow{
			Enable:    request.Body.PercentLow.Enable,
			Threshold: request.Body.PercentLow.Threshold,
		},
		BatteryHealthLow: config.BatteryHealthLow{
			Enable:    request.Body.HealthLow.Enable,
			Threshold: request.Body.HealthLow.Threshold,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("config service update battery monitoring config: %w", err)
	}

	return gen.UpdateBatteryMonitoringConfig200JSONResponse(h.convertBatteryMonitoringConfigToResponse(cfg)), nil
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
		Pic: gen.PICConfig{
			Serial:            h.convertSerialConfigToResponse(cfg.PIC.Serial),
			EnableAck:         cfg.PIC.EnableACK,
			CommandAckTimeout: int(cfg.PIC.CommandACKTimeout.Milliseconds()),
		},
		Esp: gen.ESPConfig{
			Serial:            h.convertSerialConfigToResponse(cfg.ESP.Serial),
			EnableAck:         cfg.ESP.EnableACK,
			CommandAckTimeout: int(cfg.ESP.CommandACKTimeout.Milliseconds()),
		},
		Leds: gen.LedsConfig{
			System: gen.LedConfig{
				Pin: cfg.Leds.System.Pin,
			},
			Alert: gen.LedConfig{
				Pin: cfg.Leds.Alert.Pin,
			},
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
		Enable:  cfg.Enable,
		Address: cfg.Address,
		Token:   cfg.Token,
	}
}

func (configHandler) convertHTTPConfigToResponse(cfg config.HTTP) gen.HTTPConfig {
	return gen.HTTPConfig{
		Port:    int(cfg.Port),
		Swagger: cfg.Swagger,
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
		Ip:       cfg.STA.IP,
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

func (configHandler) convertCommandConfigToResponse(cfg config.Command) gen.CommandConfig {
	return gen.CommandConfig{
		CargoLift: gen.CargoLiftConfig{
			StableReadCount: cfg.CargoLift.StableReadCount,
		},
		CargoLower: gen.CargoLowerConfig{
			StableReadCount: cfg.CargoLower.StableReadCount,
			BottomObstacleTracking: gen.ObstacleTracking{
				EnterDistance: cfg.CargoLower.BottomObstacleTracking.EnterDistance,
				ExitDistance:  cfg.CargoLower.BottomObstacleTracking.ExitDistance,
			},
		},
	}
}

func (configHandler) convertBatteryMonitoringConfigToResponse(cfg config.BatteryMonitoring) gen.BatteryMonitoringConfig {
	return gen.BatteryMonitoringConfig{
		VoltageLow: gen.BatteryVoltageLowConfig{
			Enable:    cfg.BatteryVoltageLow.Enable,
			Threshold: cfg.BatteryVoltageLow.Threshold,
		},
		VoltageHigh: gen.BatteryVoltageHighConfig{
			Enable:    cfg.BatteryVoltageHigh.Enable,
			Threshold: cfg.BatteryVoltageHigh.Threshold,
		},
		CellVoltageHigh: gen.BatteryCellVoltageHighConfig{
			Enable:    cfg.BatteryCellVoltageHigh.Enable,
			Threshold: cfg.BatteryCellVoltageHigh.Threshold,
		},
		CellVoltageLow: gen.BatteryCellVoltageLowConfig{
			Enable:    cfg.BatteryCellVoltageLow.Enable,
			Threshold: cfg.BatteryCellVoltageLow.Threshold,
		},
		CellVoltageDiff: gen.BatteryCellVoltageDiffConfig{
			Enable:    cfg.BatteryCellVoltageDiff.Enable,
			Threshold: cfg.BatteryCellVoltageDiff.Threshold,
		},
		CurrentHigh: gen.BatteryCurrentHighConfig{
			Enable:    cfg.BatteryCurrentHigh.Enable,
			Threshold: cfg.BatteryCurrentHigh.Threshold,
		},
		TempHigh: gen.BatteryTempHighConfig{
			Enable:    cfg.BatteryTempHigh.Enable,
			Threshold: cfg.BatteryTempHigh.Threshold,
		},
		PercentLow: gen.BatteryPercentLowConfig{
			Enable:    cfg.BatteryPercentLow.Enable,
			Threshold: cfg.BatteryPercentLow.Threshold,
		},
		HealthLow: gen.BatteryHealthLowConfig{
			Enable:    cfg.BatteryHealthLow.Enable,
			Threshold: cfg.BatteryHealthLow.Threshold,
		},
	}
}
