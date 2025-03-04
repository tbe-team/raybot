package serviceimpl

import (
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/pkg/validator"
)

type serviceImpl struct {
	robotService *robotService
}

func New(repo repository.Repository, validator validator.Validator) service.Service {
	return &serviceImpl{
		robotService: newRobotService(repo.RobotState(), validator),
	}
}

func (s *serviceImpl) RobotService() service.RobotService {
	return s.robotService
}
