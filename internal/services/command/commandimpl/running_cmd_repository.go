package commandimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/command"
)

type runningCmdRepository struct {
	cmd *command.CancelableCommand
	mu  sync.RWMutex
}

func NewRunningCmdRepository() command.RunningCommandRepository {
	return &runningCmdRepository{
		cmd: nil,
	}
}

func (r *runningCmdRepository) Add(_ context.Context, cmd command.CancelableCommand) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cmd != nil {
		return command.ErrRunningCommandExists
	}
	r.cmd = &cmd
	return nil
}

func (r *runningCmdRepository) Get(_ context.Context) (command.CancelableCommand, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.cmd == nil {
		return command.CancelableCommand{}, command.ErrRunningCommandNotFound
	}
	return *r.cmd, nil
}

func (r *runningCmdRepository) Remove(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cmd = nil
	return nil
}
