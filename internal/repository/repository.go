package repository

type Repository interface {
	RobotState() RobotStateRepository
	PICSerialCommand() PICSerialCommandRepository
}
