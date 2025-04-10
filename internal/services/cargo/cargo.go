package cargo

import "context"

type UpdateCargoDoorParams struct {
	IsOpen bool
}

type UpdateCargoQRCodeParams struct {
	QRCode string `validate:"required"`
}

type UpdateCargoBottomDistanceParams struct {
	BottomDistance uint16 `validate:"min=0"`
}

type UpdateCargoDoorMotorStateParams struct {
	Direction    DoorDirection `validate:"enum"`
	SetDirection bool
	Speed        uint8 `validate:"min=0,max=100"`
	SetSpeed     bool
	IsRunning    bool
	SetIsRunning bool
	Enabled      bool
	SetEnabled   bool
}

type Service interface {
	UpdateCargoDoor(ctx context.Context, params UpdateCargoDoorParams) error
	UpdateCargoQRCode(ctx context.Context, params UpdateCargoQRCodeParams) error
	UpdateCargoBottomDistance(ctx context.Context, params UpdateCargoBottomDistanceParams) error
	UpdateCargoDoorMotorState(ctx context.Context, params UpdateCargoDoorMotorStateParams) error

	OpenCargoDoor(ctx context.Context) error
	CloseCargoDoor(ctx context.Context) error
}

type Repository interface {
	GetCargo(ctx context.Context) (Cargo, error)
	GetCargoDoorMotorState(ctx context.Context) (DoorMotorState, error)
	UpdateCargoDoor(ctx context.Context, params UpdateCargoDoorParams) error
	UpdateCargoQRCode(ctx context.Context, params UpdateCargoQRCodeParams) error
	UpdateCargoBottomDistance(ctx context.Context, params UpdateCargoBottomDistanceParams) error
	UpdateCargoDoorMotorState(ctx context.Context, params UpdateCargoDoorMotorStateParams) error
}
