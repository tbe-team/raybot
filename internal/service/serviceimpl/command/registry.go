package command

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/model"
)

type executor interface {
	Execute(ctx context.Context, command model.Command) error
}

type registry struct {
	executors map[model.CommandType]executor
}

func newRegistry() *registry {
	return &registry{
		executors: make(map[model.CommandType]executor),
	}
}

func (r *registry) Register(commandType model.CommandType, executor executor) {
	r.executors[commandType] = executor
}

func (r *registry) GetExecutor(commandType model.CommandType) (executor, error) {
	executor, ok := r.executors[commandType]
	if !ok {
		return nil, fmt.Errorf("command executor not found")
	}
	return executor, nil
}
