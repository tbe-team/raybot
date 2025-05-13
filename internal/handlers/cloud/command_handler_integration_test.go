package cloud_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/tbe-team/raybot/internal/handlers/cloud/cloudtest"
	commandv1 "github.com/tbe-team/raybot/internal/handlers/cloud/gen/command/v1"
	"github.com/tbe-team/raybot/internal/services/command"
)

func TestIntegrationCommandHandler_CreateCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	testEnv := cloudtest.SetupTunnelTestEnv(t)
	client := commandv1.NewCommandServiceClient(testEnv.TunnelChannel)

	t.Run("Create command", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *commandv1.CreateCommandRequest
			wantErr bool
		}{
			{
				name: "Should create stop movement command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_STOP_MOVEMENT,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_Stop{
							Stop: &commandv1.StopInputs{},
						},
					},
				},
			},
			{
				name: "Should create move to command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_MOVE_TO,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_MoveTo{
							MoveTo: &commandv1.MoveToInputs{
								Location:  "test-location",
								Direction: commandv1.MoveToInputs_DIRECTION_FORWARD,
							},
						},
					},
				},
			},
			{
				name: "Should create move forward command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_MOVE_FORWARD,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_MoveForward{},
					},
				},
			},
			{
				name: "Should create move backward command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_MOVE_BACKWARD,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_MoveBackward{},
					},
				},
			},
			{
				name: "Should create cargo open command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_CARGO_OPEN,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_CargoOpen{},
					},
				},
			},
			{
				name: "Should create cargo close command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_CARGO_CLOSE,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_CargoClose{},
					},
				},
			},
			{
				name: "Should create cargo lift command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_CARGO_LIFT,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_CargoLift{},
					},
				},
			},
			{
				name: "Should create cargo lower command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_CARGO_LOWER,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_CargoLower{},
					},
				},
			},
			{
				name: "Should create cargo check QR command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_CARGO_CHECK_QR,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_CargoCheckQr{
							CargoCheckQr: &commandv1.CargoCheckQRInputs{
								QrCode: "test-qr-code",
							},
						},
					},
				},
			},
			{
				name: "Should create scan location command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_SCAN_LOCATION,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_ScanLocation{},
					},
				},
			},
			{
				name: "Should create wait command successfully",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_WAIT,
					Inputs: &commandv1.CommandInputs{
						Inputs: &commandv1.CommandInputs_Wait{
							Wait: &commandv1.WaitInputs{
								DurationMs: 1000,
							},
						},
					},
				},
			},
			{
				name: "Should return error for invalid command type",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_UNSPECIFIED,
				},
				wantErr: true,
			},
			{
				name: "Should return error for missing move to inputs",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_MOVE_TO,
				},
				wantErr: true,
			},
			{
				name: "Should return error for missing cargo check QR inputs",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_CARGO_CHECK_QR,
				},
				wantErr: true,
			},
			{
				name: "Should return error for missing wait inputs",
				req: &commandv1.CreateCommandRequest{
					Type: commandv1.CommandType_COMMAND_TYPE_WAIT,
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				createResp, err := client.CreateCommand(context.Background(), tt.req)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, createResp)
				require.NotEmpty(t, createResp.Command.Id)

				// Get command
				getResp, err := client.GetCommand(context.Background(), &commandv1.GetCommandRequest{
					Id: createResp.Command.Id,
				})
				require.NoError(t, err)
				require.NotNil(t, getResp)
				require.Equal(t, createResp.Command.Id, getResp.Command.Id)
				require.Equal(t, tt.req.Type, getResp.Command.Type)
				require.Equal(t, commandv1.CommandSource_COMMAND_SOURCE_CLOUD, getResp.Command.Source)
				require.True(t, proto.Equal(tt.req.Inputs, getResp.Command.Inputs))
			})
		}
	})

	t.Run("Get command", func(t *testing.T) {
		createdCommand, err := testEnv.CommandService.CreateCommand(context.Background(), command.CreateCommandParams{
			Source: command.SourceApp,
			Inputs: command.StopMovementInputs{},
		})
		require.NoError(t, err)
		require.NotEmpty(t, createdCommand)

		const nonExistentCommandID = 99999

		tests := []struct {
			name    string
			req     *commandv1.GetCommandRequest
			wantErr bool
		}{
			{
				name:    "Should get command successfully",
				req:     &commandv1.GetCommandRequest{Id: createdCommand.ID},
				wantErr: false,
			},
			{
				name:    "Should return error for non-existent command",
				req:     &commandv1.GetCommandRequest{Id: nonExistentCommandID},
				wantErr: true,
			},
			{
				name:    "Should return error for missing command ID",
				req:     &commandv1.GetCommandRequest{},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				getResp, err := client.GetCommand(context.Background(), tt.req)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, getResp)
				require.NotEmpty(t, getResp.Command)
				require.Equal(t, createdCommand.ID, getResp.Command.Id)
			})
		}
	})
}
