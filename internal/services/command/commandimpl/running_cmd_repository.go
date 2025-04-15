package commandimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/services/command"
)

type runningCommand struct {
	command.Command

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newRunningCommand(cmd command.Command) *runningCommand {
	ctx, cancel := context.WithCancel(context.Background())
	return &runningCommand{
		Command:    cmd,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Cancel cancels the running command.
func (c *runningCommand) Cancel() {
	c.cancelFunc()
}

// Context returns the context of the running command.
func (c *runningCommand) Context() context.Context {
	return c.ctx
}

type runningCmdRepository struct {
	cmd *runningCommand
	mu  sync.RWMutex
}

func newRunningCmdRepository() *runningCmdRepository {
	return &runningCmdRepository{
		cmd: nil,
	}
}

func (r *runningCmdRepository) Add(cmd *runningCommand) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cmd = cmd
}

func (r *runningCmdRepository) Get() *runningCommand {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.cmd
}

func (r *runningCmdRepository) Remove() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cmd = nil
}
