package wifi

import (
	"context"
)

type Service interface {
	// Run initializes the wifi connection.
	// This will set the wifi to AP mode if the config is set to AP mode.
	// Otherwise, it will set the wifi to STA mode.
	Run(ctx context.Context) error
}
