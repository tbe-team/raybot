package repoimpl

import (
	"context"
	"sync"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var (
	ErrESPSerialCommandNotFound = xerror.NotFound(nil, "espSerialCommand.notFound", "ESP serial command not found")
)

type ESPSerialCommandRepository struct {
	mu     sync.RWMutex
	cmdMap map[string]model.ESPSerialCommand
}

func NewESPSerialCommandRepository() *ESPSerialCommandRepository {
	return &ESPSerialCommandRepository{cmdMap: make(map[string]model.ESPSerialCommand)}
}

func (r *ESPSerialCommandRepository) GetESPSerialCommand(_ context.Context, id string) (model.ESPSerialCommand, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, ok := r.cmdMap[id]
	if !ok {
		return model.ESPSerialCommand{}, ErrESPSerialCommandNotFound
	}

	return cmd, nil
}

func (r *ESPSerialCommandRepository) CreateESPSerialCommand(_ context.Context, command model.ESPSerialCommand) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cmdMap[command.ID] = command

	return nil
}

func (r *ESPSerialCommandRepository) DeleteESPSerialCommand(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.cmdMap, id)

	return nil
}
