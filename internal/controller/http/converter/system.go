package converter

import (
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/service"
)

func ToSystemConfigResponse(cfg service.GetSystemConfigOutput) gen.SystemConfigResponse {
	return gen.SystemConfigResponse{
		Grpc: gen.GRPCConfig{
			Port: cfg.GRPCConfig.Port,
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
	}
}
