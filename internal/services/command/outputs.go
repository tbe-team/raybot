package command

var (
	_ Outputs = (*StopMovementOutputs)(nil)
	_ Outputs = (*MoveForwardOutputs)(nil)
	_ Outputs = (*MoveBackwardOutputs)(nil)
	_ Outputs = (*MoveToOutputs)(nil)
	_ Outputs = (*CargoOpenOutputs)(nil)
	_ Outputs = (*CargoCloseOutputs)(nil)
	_ Outputs = (*CargoLiftOutputs)(nil)
	_ Outputs = (*CargoLowerOutputs)(nil)
	_ Outputs = (*CargoCheckQROutputs)(nil)
)

type Outputs interface {
	isOutputs()
	CommandType() CommandType
}

type StopMovementOutputs struct{}

func (StopMovementOutputs) CommandType() CommandType {
	return CommandTypeStopMovement
}

func (StopMovementOutputs) isOutputs() {}

type MoveForwardOutputs struct{}

func (MoveForwardOutputs) CommandType() CommandType {
	return CommandTypeMoveForward
}

func (MoveForwardOutputs) isOutputs() {}

type MoveBackwardOutputs struct{}

func (MoveBackwardOutputs) CommandType() CommandType {
	return CommandTypeMoveBackward
}

func (MoveBackwardOutputs) isOutputs() {}

type MoveToOutputs struct{}

func (MoveToOutputs) CommandType() CommandType {
	return CommandTypeMoveTo
}

func (MoveToOutputs) isOutputs() {}

type CargoOpenOutputs struct{}

func (CargoOpenOutputs) CommandType() CommandType {
	return CommandTypeCargoOpen
}

func (CargoOpenOutputs) isOutputs() {}

type CargoCloseOutputs struct{}

func (CargoCloseOutputs) CommandType() CommandType {
	return CommandTypeCargoClose
}

func (CargoCloseOutputs) isOutputs() {}

type CargoLiftOutputs struct{}

func (CargoLiftOutputs) CommandType() CommandType {
	return CommandTypeCargoLift
}

func (CargoLiftOutputs) isOutputs() {}

type CargoLowerOutputs struct{}

func (CargoLowerOutputs) CommandType() CommandType {
	return CommandTypeCargoLower
}

func (CargoLowerOutputs) isOutputs() {}

type CargoCheckQROutputs struct{}

func (CargoCheckQROutputs) CommandType() CommandType {
	return CommandTypeCargoCheckQR
}

func (CargoCheckQROutputs) isOutputs() {}

func UnmarshalOutputs(cmdType CommandType, outputsBytes []byte) (Outputs, error) {
	var outputs Outputs

	switch cmdType {
	case CommandTypeStopMovement:
		outputs = &StopMovementOutputs{}

	case CommandTypeMoveForward:
		outputs = &MoveForwardOutputs{}

	case CommandTypeMoveBackward:
		outputs = &MoveBackwardOutputs{}

	case CommandTypeMoveTo:
		outputs = &MoveToOutputs{}

	case CommandTypeCargoOpen:
		outputs = &CargoOpenOutputs{}

	case CommandTypeCargoClose:
		outputs = &CargoCloseOutputs{}

	case CommandTypeCargoLift:
		outputs = &CargoLiftOutputs{}

	case CommandTypeCargoLower:
		outputs = &CargoLowerOutputs{}

	case CommandTypeCargoCheckQR:
		outputs = &CargoCheckQROutputs{}
	}

	return outputs, nil
}
