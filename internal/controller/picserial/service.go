package pic

import (
	"flag"

	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
)

type Config struct {
	Serial serial.Config `yaml:"serial"`
}

// RegisterFlags registers flags for the PIC configuration.
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	cfg.Serial.RegisterFlagsWithPrefix("pic.", f)
}

// Validate validates the PIC configuration.
func (cfg *Config) Validate() error {
	return cfg.Serial.Validate()
}

type PICSerialService struct {
	cfg Config

	serialClient serial.Client
}
