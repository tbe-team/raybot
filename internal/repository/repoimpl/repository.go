package repoimpl

import (
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type repo struct {
	robotStateRepo       repository.RobotStateRepository
	picSerialCommandRepo repository.PICSerialCommandRepository
}

func New() repository.Repository {
	queries := sqlc.New()
	return &repo{
		robotStateRepo:       NewRobotStateRepository(queries),
		picSerialCommandRepo: NewPICSerialCommandRepository(),
	}
}

func (r *repo) RobotState() repository.RobotStateRepository {
	return r.robotStateRepo
}

func (r *repo) PICSerialCommand() repository.PICSerialCommandRepository {
	return r.picSerialCommandRepo
}
