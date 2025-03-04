package repoimpl

import "github.com/tbe-team/raybot/internal/repository"

type repo struct {
	robotStateRepo repository.RobotStateRepository
}

func New() repository.Repository {
	return &repo{
		robotStateRepo: NewRobotStateRepository(),
	}
}

func (r *repo) RobotState() repository.RobotStateRepository {
	return r.robotStateRepo
}
