package cloud

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	commandv1 "github.com/tbe-team/raybot/internal/handlers/cloud/gen/command/v1"
	"github.com/tbe-team/raybot/internal/services/command"
)

type commandHandler struct {
	commandv1.UnimplementedCommandServiceServer
	commandService command.Service
}

func newCommandHandler(commandService command.Service) commandv1.CommandServiceServer {
	return &commandHandler{
		commandService: commandService,
	}
}

func (h commandHandler) CreateCommand(ctx context.Context, req *commandv1.CreateCommandRequest) (*commandv1.CreateCommandResponse, error) {
	inputs, err := h.convertReqInputsToCommandInputs(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "convert req inputs to command inputs: %v", err)
	}
	cmd, err := h.commandService.CreateCommand(ctx, command.CreateCommandParams{
		Source: command.SourceCloud,
		Inputs: inputs,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create command: %v", err)
	}
	return &commandv1.CreateCommandResponse{
		Command: h.convertCommandToResponse(cmd),
	}, nil
}

func (h commandHandler) GetCommand(ctx context.Context, req *commandv1.GetCommandRequest) (*commandv1.GetCommandResponse, error) {
	cmd, err := h.commandService.GetCommandByID(ctx, command.GetCommandByIDParams{
		CommandID: req.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get command: %v", err)
	}
	return &commandv1.GetCommandResponse{
		Command: h.convertCommandToResponse(cmd),
	}, nil
}

func (commandHandler) convertReqInputsToCommandInputs(req *commandv1.CreateCommandRequest) (command.Inputs, error) {
	switch req.Type {
	case commandv1.CommandType_COMMAND_TYPE_STOP_MOVEMENT:
		return &command.StopMovementInputs{}, nil

	case commandv1.CommandType_COMMAND_TYPE_MOVE_TO:
		i := req.Inputs.GetMoveTo()
		if i == nil {
			return nil, fmt.Errorf("move to inputs is nil")
		}
		return &command.MoveToInputs{
			Location: i.Location,
		}, nil

	case commandv1.CommandType_COMMAND_TYPE_MOVE_FORWARD:
		return &command.MoveForwardInputs{}, nil

	case commandv1.CommandType_COMMAND_TYPE_MOVE_BACKWARD:
		return &command.MoveBackwardInputs{}, nil

	case commandv1.CommandType_COMMAND_TYPE_CARGO_OPEN:
		return &command.CargoOpenInputs{}, nil

	case commandv1.CommandType_COMMAND_TYPE_CARGO_CLOSE:
		return &command.CargoCloseInputs{}, nil

	case commandv1.CommandType_COMMAND_TYPE_CARGO_LIFT:
		return &command.CargoLiftInputs{}, nil

	case commandv1.CommandType_COMMAND_TYPE_CARGO_LOWER:
		return &command.CargoLowerInputs{}, nil

	case commandv1.CommandType_COMMAND_TYPE_CARGO_CHECK_QR:
		i := req.Inputs.GetCargoCheckQr()
		if i == nil {
			return nil, fmt.Errorf("cargo check qr inputs is nil")
		}
		return &command.CargoCheckQRInputs{
			QRCode: i.QrCode,
		}, nil

	case commandv1.CommandType_COMMAND_TYPE_SCAN_LOCATION:
		return &command.ScanLocationInputs{}, nil

	case commandv1.CommandType_COMMAND_TYPE_WAIT:
		i := req.Inputs.GetWait()
		if i == nil {
			return nil, fmt.Errorf("wait inputs is nil")
		}
		return &command.WaitInputs{
			DurationMs: i.DurationMs,
		}, nil

	default:
		return nil, fmt.Errorf("invalid command type: %v", req.Type)
	}
}

func (h commandHandler) convertCommandToResponse(cmd command.Command) *commandv1.Command {
	var (
		startedAt   *timestamppb.Timestamp
		completedAt *timestamppb.Timestamp
	)
	if cmd.StartedAt != nil {
		startedAt = timestamppb.New(*cmd.StartedAt)
	}
	if cmd.CompletedAt != nil {
		completedAt = timestamppb.New(*cmd.CompletedAt)
	}

	return &commandv1.Command{
		Id:          cmd.ID,
		Type:        h.convertCommandTypeToResponse(cmd.Type),
		Inputs:      h.convertCommandInputsToResponse(cmd.Inputs),
		Status:      h.convertCommandStatusToResponse(cmd.Status),
		Source:      h.convertCommandSourceToResponse(cmd.Source),
		Outputs:     h.convertCommandOutputsToResponse(cmd.Outputs),
		Error:       cmd.Error,
		StartedAt:   startedAt,
		CompletedAt: completedAt,
		CreatedAt:   timestamppb.New(cmd.CreatedAt),
		UpdatedAt:   timestamppb.New(cmd.UpdatedAt),
	}
}

func (commandHandler) convertCommandTypeToResponse(cmdType command.CommandType) commandv1.CommandType {
	switch cmdType {
	case command.CommandTypeStopMovement:
		return commandv1.CommandType_COMMAND_TYPE_STOP_MOVEMENT

	case command.CommandTypeMoveTo:
		return commandv1.CommandType_COMMAND_TYPE_MOVE_TO

	case command.CommandTypeMoveForward:
		return commandv1.CommandType_COMMAND_TYPE_MOVE_FORWARD

	case command.CommandTypeMoveBackward:
		return commandv1.CommandType_COMMAND_TYPE_MOVE_BACKWARD

	case command.CommandTypeCargoOpen:
		return commandv1.CommandType_COMMAND_TYPE_CARGO_OPEN

	case command.CommandTypeCargoClose:
		return commandv1.CommandType_COMMAND_TYPE_CARGO_CLOSE

	case command.CommandTypeCargoLift:
		return commandv1.CommandType_COMMAND_TYPE_CARGO_LIFT

	case command.CommandTypeCargoLower:
		return commandv1.CommandType_COMMAND_TYPE_CARGO_LOWER

	case command.CommandTypeCargoCheckQR:
		return commandv1.CommandType_COMMAND_TYPE_CARGO_CHECK_QR

	case command.CommandTypeScanLocation:
		return commandv1.CommandType_COMMAND_TYPE_SCAN_LOCATION

	case command.CommandTypeWait:
		return commandv1.CommandType_COMMAND_TYPE_WAIT

	default:
		return commandv1.CommandType_COMMAND_TYPE_UNSPECIFIED
	}
}

func (commandHandler) convertCommandInputsToResponse(inputs command.Inputs) *commandv1.CommandInputs {
	switch i := inputs.(type) {
	case *command.StopMovementInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_Stop{
				Stop: &commandv1.StopInputs{},
			},
		}

	case *command.MoveForwardInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_MoveForward{
				MoveForward: &commandv1.MoveForwardInputs{},
			},
		}

	case *command.MoveBackwardInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_MoveBackward{
				MoveBackward: &commandv1.MoveBackwardInputs{},
			},
		}

	case *command.MoveToInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_MoveTo{
				MoveTo: &commandv1.MoveToInputs{
					Location: i.Location,
				},
			},
		}

	case *command.CargoOpenInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_CargoOpen{
				CargoOpen: &commandv1.CargoOpenInputs{},
			},
		}

	case *command.CargoCloseInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_CargoClose{
				CargoClose: &commandv1.CargoCloseInputs{},
			},
		}

	case *command.CargoLiftInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_CargoLift{
				CargoLift: &commandv1.CargoLiftInputs{},
			},
		}

	case *command.CargoLowerInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_CargoLower{
				CargoLower: &commandv1.CargoLowerInputs{},
			},
		}

	case *command.CargoCheckQRInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_CargoCheckQr{
				CargoCheckQr: &commandv1.CargoCheckQRInputs{
					QrCode: i.QRCode,
				},
			},
		}

	case *command.ScanLocationInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_ScanLocation{
				ScanLocation: &commandv1.ScanLocationInputs{},
			},
		}

	case *command.WaitInputs:
		return &commandv1.CommandInputs{
			Inputs: &commandv1.CommandInputs_Wait{
				Wait: &commandv1.WaitInputs{
					DurationMs: i.DurationMs,
				},
			},
		}

	default:
		return nil
	}
}

func (commandHandler) convertCommandStatusToResponse(status command.Status) commandv1.CommandStatus {
	switch status {
	case command.StatusQueued:
		return commandv1.CommandStatus_COMMAND_STATUS_QUEUED

	case command.StatusProcessing:
		return commandv1.CommandStatus_COMMAND_STATUS_PROCESSING

	case command.StatusSucceeded:
		return commandv1.CommandStatus_COMMAND_STATUS_SUCCEEDED

	case command.StatusFailed:
		return commandv1.CommandStatus_COMMAND_STATUS_FAILED

	default:
		return commandv1.CommandStatus_COMMAND_STATUS_UNSPECIFIED
	}
}

func (commandHandler) convertCommandSourceToResponse(source command.Source) commandv1.CommandSource {
	switch source {
	case command.SourceCloud:
		return commandv1.CommandSource_COMMAND_SOURCE_CLOUD

	case command.SourceApp:
		return commandv1.CommandSource_COMMAND_SOURCE_APP

	default:
		return commandv1.CommandSource_COMMAND_SOURCE_UNSPECIFIED
	}
}

func (commandHandler) convertCommandOutputsToResponse(outputs command.Outputs) *commandv1.CommandOutputs {
	switch o := outputs.(type) {
	case *command.StopMovementOutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_Stop{
				Stop: &commandv1.StopOutputs{},
			},
		}

	case *command.MoveForwardOutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_MoveForward{
				MoveForward: &commandv1.MoveForwardOutputs{},
			},
		}

	case *command.MoveBackwardOutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_MoveBackward{
				MoveBackward: &commandv1.MoveBackwardOutputs{},
			},
		}

	case *command.MoveToOutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_MoveTo{
				MoveTo: &commandv1.MoveToOutputs{},
			},
		}

	case *command.CargoOpenOutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_CargoOpen{
				CargoOpen: &commandv1.CargoOpenOutputs{},
			},
		}

	case *command.CargoCloseOutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_CargoClose{
				CargoClose: &commandv1.CargoCloseOutputs{},
			},
		}

	case *command.CargoLiftOutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_CargoLift{
				CargoLift: &commandv1.CargoLiftOutputs{},
			},
		}

	case *command.CargoLowerOutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_CargoLower{
				CargoLower: &commandv1.CargoLowerOutputs{},
			},
		}

	case *command.CargoCheckQROutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_CargoCheckQr{
				CargoCheckQr: &commandv1.CargoCheckQROutputs{},
			},
		}

	case *command.ScanLocationOutputs:
		locs := []*commandv1.Location{}
		for _, l := range o.Locations {
			locs = append(locs, &commandv1.Location{
				Location:  l.Location,
				ScannedAt: timestamppb.New(l.ScannedAt),
			})
		}

		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_ScanLocation{
				ScanLocation: &commandv1.ScanLocationOutputs{
					Locations: locs,
				},
			},
		}

	case *command.WaitOutputs:
		return &commandv1.CommandOutputs{
			Outputs: &commandv1.CommandOutputs_Wait{
				Wait: &commandv1.WaitOutputs{},
			},
		}

	default:
		return nil
	}
}
