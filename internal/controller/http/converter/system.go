package converter

import (
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/service"
)

func ConvertSystemConfigToResponse(cfg service.GetSystemConfigOutput) gen.SystemConfigResponse {
	return gen.SystemConfigResponse{
		Grpc: gen.GRPCConfig{
			Server: gen.GRPCServerConfig{
				Enable: cfg.GRPCConfig.Server.Enable,
			},
			Cloud: gen.CloudConfig{
				Address: cfg.GRPCConfig.Cloud.Address,
				Token:   cfg.GRPCConfig.Cloud.Token,
			},
		},
		Http: gen.HTTPConfig{
			EnableSwagger: cfg.HTTPConfig.EnableSwagger,
		},
		Log: gen.LogConfig{
			Level:     cfg.LogConfig.Level,
			Format:    cfg.LogConfig.Format,
			AddSource: cfg.LogConfig.AddSource,
		},
		Pic: gen.PicConfig{
			Serial: gen.SerialConfig{
				Port:        cfg.PICConfig.Serial.Port,
				BaudRate:    cfg.PICConfig.Serial.BaudRate,
				DataBits:    cfg.PICConfig.Serial.DataBits,
				StopBits:    cfg.PICConfig.Serial.StopBits,
				Parity:      cfg.PICConfig.Serial.Parity,
				ReadTimeout: cfg.PICConfig.Serial.ReadTimeout.Seconds(),
			},
		},
		Esp: gen.ESPConfig{
			Serial: gen.SerialConfig{
				Port:        cfg.ESPConfig.Serial.Port,
				BaudRate:    cfg.ESPConfig.Serial.BaudRate,
				DataBits:    cfg.ESPConfig.Serial.DataBits,
				StopBits:    cfg.ESPConfig.Serial.StopBits,
				Parity:      cfg.ESPConfig.Serial.Parity,
				ReadTimeout: cfg.ESPConfig.Serial.ReadTimeout.Seconds(),
			},
		},
	}
}
