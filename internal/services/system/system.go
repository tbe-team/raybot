package system

import "context"

type Service interface {
	RestartApplication(ctx context.Context) error
}
