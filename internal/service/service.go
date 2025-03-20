package service

type Service interface {
	RobotService() RobotService
	SystemService() SystemService
	PICService() PICService
}
