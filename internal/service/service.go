package service

type Service interface {
	RobotStateService() RobotStateService
	SystemService() SystemService
	PICService() PICService
	LocationService() LocationService
	CommandService() CommandService
	CargoControlService() CargoControlService
}
