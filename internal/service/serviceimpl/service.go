package serviceimpl

import (
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/pkg/validator"
)

type serviceImpl struct {
	robotService *RobotService
}

func New(repo repository.Repository, validator validator.Validator) service.Service {
	return &serviceImpl{
		robotService: NewRobotService(repo.RobotState(), validator),
	}
}

func (s *serviceImpl) RobotService() service.RobotService {
	return s.robotService
}
