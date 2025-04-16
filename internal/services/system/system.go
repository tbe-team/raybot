package system

import "context"

type Service interface {
	Reboot(ctx context.Context) error
}
