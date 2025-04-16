package systemimpl

import (
	"context"
	"log/slog"
	"os/exec"
	"time"

	"github.com/tbe-team/raybot/internal/services/system"
)

type service struct {
	log *slog.Logger
}

func NewService(log *slog.Logger) system.Service {
	return &service{log: log}
}

func (s service) Reboot(_ context.Context) error {
	go func() {
		time.Sleep(1 * time.Second)
		cmd := exec.Command("reboot")
		if err := cmd.Run(); err != nil {
			s.log.Error("failed to reboot", slog.Any("error", err))
		}
	}()

	return nil
}
