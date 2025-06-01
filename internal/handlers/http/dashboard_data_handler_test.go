package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/appstate"
	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/dashboarddata"
	dashboarddatamocks "github.com/tbe-team/raybot/internal/services/dashboarddata/mocks"
	"github.com/tbe-team/raybot/internal/services/distancesensor"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/internal/services/location"
	"github.com/tbe-team/raybot/pkg/ptr"
)

func TestDashboardDataHandler_GetRobotState(t *testing.T) {
	t.Run("Should get robot state successfully", func(t *testing.T) {
		dashboardDataService := dashboarddatamocks.NewFakeService(t)
		dashboardDataService.EXPECT().GetRobotState(mock.Anything).Return(validRobotState, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.dashboardDataService = dashboardDataService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/robot-state", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.GetRobotState200JSONResponse](t, rec.Body)
		require.Equal(t, validRobotState.Battery.Current, res.Battery.Current)
		require.Equal(t, validRobotState.Battery.Voltage, res.Battery.Voltage)
		require.Equal(t, validRobotState.Battery.Temp, res.Battery.Temp)
		require.Equal(t, validRobotState.Battery.CellVoltages, res.Battery.CellVoltages)
		require.Equal(t, validRobotState.Battery.Percent, res.Battery.Percent)
		require.Equal(t, validRobotState.Battery.Fault, res.Battery.Fault)
		require.Equal(t, validRobotState.Battery.Health, res.Battery.Health)
		require.NotEmpty(t, validRobotState.Battery.UpdatedAt)

		require.Equal(t, validRobotState.BatteryCharge.CurrentLimit, res.Charge.CurrentLimit)
		require.Equal(t, validRobotState.BatteryCharge.Enabled, res.Charge.Enabled)
		require.NotEmpty(t, validRobotState.BatteryCharge.UpdatedAt)

		require.Equal(t, validRobotState.BatteryDischarge.CurrentLimit, res.Discharge.CurrentLimit)
		require.Equal(t, validRobotState.BatteryDischarge.Enabled, res.Discharge.Enabled)
		require.NotEmpty(t, validRobotState.BatteryDischarge.UpdatedAt)

		require.Equal(t, validRobotState.DistanceSensor.FrontDistance, res.DistanceSensor.FrontDistance)
		require.Equal(t, validRobotState.DistanceSensor.BackDistance, res.DistanceSensor.BackDistance)
		require.Equal(t, validRobotState.DistanceSensor.DownDistance, res.DistanceSensor.DownDistance)
		require.NotEmpty(t, validRobotState.DistanceSensor.UpdatedAt)

		require.Equal(t, validRobotState.LiftMotor.CurrentPosition, res.LiftMotor.CurrentPosition)
		require.Equal(t, validRobotState.LiftMotor.TargetPosition, res.LiftMotor.TargetPosition)
		require.Equal(t, validRobotState.LiftMotor.IsRunning, res.LiftMotor.IsRunning)
		require.Equal(t, validRobotState.LiftMotor.Enabled, res.LiftMotor.Enabled)
		require.NotEmpty(t, validRobotState.LiftMotor.UpdatedAt)

		require.Equal(t, validRobotState.DriveMotor.Direction.String(), "BACKWARD")
		require.Equal(t, validRobotState.DriveMotor.Speed, res.DriveMotor.Speed)
		require.Equal(t, validRobotState.DriveMotor.IsRunning, res.DriveMotor.IsRunning)
		require.Equal(t, validRobotState.DriveMotor.Enabled, res.DriveMotor.Enabled)
		require.NotEmpty(t, validRobotState.DriveMotor.UpdatedAt)

		require.Equal(t, validRobotState.Location.CurrentLocation, res.Location.CurrentLocation)
		require.NotEmpty(t, validRobotState.Location.UpdatedAt)

		require.Equal(t, validRobotState.Cargo.IsOpen, res.Cargo.IsOpen)
		require.Equal(t, validRobotState.Cargo.QRCode, res.Cargo.QrCode)
		require.Equal(t, validRobotState.Cargo.BottomDistance, res.Cargo.BottomDistance)
		require.NotEmpty(t, validRobotState.Cargo.UpdatedAt)

		require.Equal(t, validRobotState.CargoDoorMotor.Direction.String(), "CLOSE")
		require.Equal(t, validRobotState.CargoDoorMotor.Speed, res.CargoDoorMotor.Speed)
		require.Equal(t, validRobotState.CargoDoorMotor.IsRunning, res.CargoDoorMotor.IsRunning)
		require.Equal(t, validRobotState.CargoDoorMotor.Enabled, res.CargoDoorMotor.Enabled)
		require.NotEmpty(t, validRobotState.CargoDoorMotor.UpdatedAt)

		require.Equal(t, validRobotState.AppState.CloudConnection.Connected, res.AppConnection.CloudConnection.Connected)
		require.NotNil(t, validRobotState.AppState.CloudConnection.LastConnectedAt)
		require.NotNil(t, validRobotState.AppState.CloudConnection.Error)

		require.Equal(t, validRobotState.AppState.ESPSerialConnection.Connected, res.AppConnection.EspSerialConnection.Connected)
		require.NotNil(t, validRobotState.AppState.ESPSerialConnection.LastConnectedAt)
		require.NotNil(t, validRobotState.AppState.ESPSerialConnection.Error)

		require.Equal(t, validRobotState.AppState.PICSerialConnection.Connected, res.AppConnection.PicSerialConnection.Connected)
		require.NotNil(t, validRobotState.AppState.PICSerialConnection.LastConnectedAt)
		require.NotNil(t, validRobotState.AppState.PICSerialConnection.Error)

		require.Equal(t, validRobotState.AppState.RFIDUSBConnection.Connected, res.AppConnection.RfidUsbConnection.Connected)
		require.NotNil(t, validRobotState.AppState.RFIDUSBConnection.LastConnectedAt)
		require.NotNil(t, validRobotState.AppState.RFIDUSBConnection.Error)
	})

	t.Run("Should not able to get robot state if fetching failed", func(t *testing.T) {
		dashboardDataService := dashboarddatamocks.NewFakeService(t)
		dashboardDataService.EXPECT().GetRobotState(mock.Anything).
			Return(dashboarddata.RobotState{}, errors.New("fetching failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.dashboardDataService = dashboardDataService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/robot-state", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestDashboardDataHandler_getUptime(t *testing.T) {
	handler := dashboardDataHandler{}
	testCases := []struct {
		name            string
		connected       bool
		lastConnectedAt *time.Time
		expected        float32
	}{
		{
			name:            "Should return 0 if not connected",
			connected:       false,
			lastConnectedAt: nil,
			expected:        0,
		},
		{
			name:            "Should return 0 if connected but last connected at nil",
			connected:       true,
			lastConnectedAt: nil,
			expected:        0,
		},
		{
			name:            "Should return uptime if connected",
			connected:       true,
			lastConnectedAt: ptr.New(time.Now().Add(-time.Hour * 24)),
			expected:        float32(time.Hour.Seconds() * 24),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uptime := handler.getUptime(tc.connected, tc.lastConnectedAt)
			require.Equal(t, tc.expected, uptime)
		})
	}
}

var validRobotState = dashboarddata.RobotState{
	Battery: battery.BatteryState{
		Current:      100,
		Voltage:      100,
		Temp:         25,
		CellVoltages: []uint16{50, 60, 70, 80},
		Percent:      100,
		Fault:        0,
		Health:       100,
		UpdatedAt:    time.Now(),
	},
	BatteryCharge: battery.ChargeSetting{
		CurrentLimit: 100,
		Enabled:      true,
		UpdatedAt:    time.Now(),
	},
	BatteryDischarge: battery.DischargeSetting{
		CurrentLimit: 100,
		Enabled:      true,
		UpdatedAt:    time.Now(),
	},
	DistanceSensor: distancesensor.DistanceSensorState{
		FrontDistance: 50,
		BackDistance:  60,
		DownDistance:  70,
		UpdatedAt:     time.Now(),
	},
	LiftMotor: liftmotor.LiftMotorState{
		CurrentPosition: 100,
		TargetPosition:  50,
		IsRunning:       true,
		Enabled:         true,
		UpdatedAt:       time.Now(),
	},
	DriveMotor: drivemotor.DriveMotorState{
		Direction: drivemotor.DirectionBackward,
		Speed:     50,
		IsRunning: true,
		Enabled:   true,
		UpdatedAt: time.Now(),
	},
	Location: location.Location{
		CurrentLocation: "123",
		UpdatedAt:       time.Now(),
	},
	Cargo: cargo.Cargo{
		IsOpen:         true,
		QRCode:         "123",
		BottomDistance: 30,
		UpdatedAt:      time.Now(),
	},
	CargoDoorMotor: cargo.DoorMotorState{
		Direction: cargo.DirectionClose,
		Speed:     50,
		IsRunning: true,
		Enabled:   true,
		UpdatedAt: time.Now(),
	},
	AppState: appstate.AppState{
		CloudConnection: appstate.CloudConnection{
			Connected:       true,
			LastConnectedAt: ptr.New(time.Now()),
			Error:           ptr.New("error"),
		},
		ESPSerialConnection: appstate.ESPSerialConnection{
			Connected:       true,
			LastConnectedAt: ptr.New(time.Now()),
			Error:           ptr.New("error"),
		},
		PICSerialConnection: appstate.PICSerialConnection{
			Connected:       true,
			LastConnectedAt: ptr.New(time.Now()),
			Error:           ptr.New("error"),
		},
		RFIDUSBConnection: appstate.RFIDUSBConnection{
			Connected:       true,
			LastConnectedAt: ptr.New(time.Now()),
			Error:           ptr.New("error"),
		},
	},
}
