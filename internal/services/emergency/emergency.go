package emergency

import "context"

type Service interface {
	// GetEmergencyState returns the current emergency state of the robot.
	GetEmergencyState(ctx context.Context) (State, error)

	// StopOperation stops the operation of the robot.
	// It should stop all motors and other components, command processing, etc.
	StopOperation(ctx context.Context) error

	// ResumeOperation resumes the operation of the robot.
	ResumeOperation(ctx context.Context) error
}

type Repository interface {
	GetEmergencyState(ctx context.Context) (State, error)
	UpdateEmergencyState(ctx context.Context, state State) error
}
