package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service/mocks"
)

func TestRobotStateHandler_GetRobotState(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		mockSetup     func(*mocks.FakeRobotService)
		expectedError bool
	}{
		{
			name: "successful get robot state",
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().GetRobotState(ctx).Return(model.RobotState{
					Battery: model.BatteryState{
						Current:      100,
						Temp:         25,
						Voltage:      12000,
						CellVoltages: []uint16{4000, 4000, 4000},
						Percent:      80,
						Health:       100,
					},
					DriveMotor: model.DriveMotorState{
						Direction: model.DriveMotorDirectionForward,
						Speed:     50,
						IsRunning: true,
						Enabled:   true,
					},
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "robot service returns error",
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, errors.New("service error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			robotService := mocks.NewFakeRobotService(t)
			h := robotStateHandler{robotService: robotService}

			tc.mockSetup(robotService)

			resp, err := h.GetRobotState(ctx, gen.GetRobotStateRequestObject{})

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				// Type assertion to check the response type
				response, ok := resp.(gen.GetRobotState200JSONResponse)
				assert.True(t, ok)

				// Verify some fields from the response
				assert.Equal(t, uint16(100), response.Battery.Current)
				assert.Equal(t, uint8(80), response.Battery.Percent)
				assert.Equal(t, uint16(12000), response.Battery.Voltage)
				assert.Equal(t, uint8(25), response.Battery.Temp)

				assert.Equal(t, uint8(50), response.DriveMotor.Speed)
				assert.Equal(t, true, response.DriveMotor.IsRunning)
				assert.Equal(t, true, response.DriveMotor.Enabled)
				assert.Equal(t, "FORWARD", response.DriveMotor.Direction)
			}
		})
	}
}
