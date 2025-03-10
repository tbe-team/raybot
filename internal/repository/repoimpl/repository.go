package repoimpl

import "github.com/tbe-team/raybot/internal/repository"

type repo struct {
	robotStateRepo       repository.RobotStateRepository
	picSerialCommandRepo repository.PICSerialCommandRepository
}

func New() repository.Repository {
	return &repo{
		robotStateRepo:       NewRobotStateRepository(),
		picSerialCommandRepo: NewPICSerialCommandRepository(),
	}
}

func (r *repo) RobotState() repository.RobotStateRepository {
	return r.robotStateRepo
}

func (r *repo) PICSerialCommand() repository.PICSerialCommandRepository {
	return r.picSerialCommandRepo
}
