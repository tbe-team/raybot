package systemimpl

import (
	"context"
	"log/slog"
	"os"
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

func (s service) RestartApplication(_ context.Context) error {
	go func() {
		time.Sleep(3 * time.Second)
		self, err := os.Executable()
		if err != nil {
			s.log.Error("failed to get executable", slog.Any("error", err))
			return
		}

		cmd := exec.Command(self, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = os.Environ()
		if err := cmd.Start(); err != nil {
			s.log.Error("failed to restart application", slog.Any("error", err))
			return
		}

		os.Exit(0)
	}()

	return nil
}
