package repoimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var (
	ErrPICSerialCommandNotFound = xerror.NotFound(nil, "picSerialCommand.notFound", "PIC serial command not found")
)

type PICSerialCommandRepository struct {
	mu     sync.RWMutex
	cmdMap map[string]model.PICSerialCommand
}

func NewPICSerialCommandRepository() *PICSerialCommandRepository {
	return &PICSerialCommandRepository{
		cmdMap: make(map[string]model.PICSerialCommand),
	}
}

func (r *PICSerialCommandRepository) GetPICSerialCommand(_ context.Context, id string) (model.PICSerialCommand, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, ok := r.cmdMap[id]
	if !ok {
		return model.PICSerialCommand{}, ErrPICSerialCommandNotFound
	}

	return cmd, nil
}

func (r *PICSerialCommandRepository) CreatePICSerialCommand(_ context.Context, command model.PICSerialCommand) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cmdMap[command.ID] = command
	return nil
}

func (r *PICSerialCommandRepository) DeletePICSerialCommand(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.cmdMap, id)
	return nil
}
