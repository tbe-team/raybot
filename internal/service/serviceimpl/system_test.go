package serviceimpl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/config/mocks"
	"github.com/tbe-team/raybot/internal/controller/http"
	"github.com/tbe-team/raybot/internal/controller/picserial"
	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/pkg/log"
)

func TestSystemService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("test GetSystemConfig", func(t *testing.T) {
		tests := []struct {
			name          string
			mock          func(_ *mocks.FakeManager)
			expected      service.GetSystemConfigOutput
			expectedError bool
		}{
			{
				name: "successful get config",
				mock: func(cfgManager *mocks.FakeManager) {
					cfgManager.EXPECT().GetConfig().Return(config.Config{
						Log: log.Config{
							Level:     "info",
							Format:    "json",
							AddSource: true,
						},
						GRPC: config.GRPCConfig{
							Server: config.GRPCServerConfig{
								Enable: true,
							},
						},
						HTTP: http.Config{
							EnableSwagger: true,
						},
						PIC: picserial.Config{
							Serial: serial.Config{
								Port:        "/dev/tty.usbserial",
								BaudRate:    9600,
								DataBits:    8,
								StopBits:    1.0,
								Parity:      "none",
								ReadTimeout: time.Second * 1,
							},
						},
					})
				},
				expected: service.GetSystemConfigOutput{
					LogConfig: service.LogConfig{
						Level:     "info",
						Format:    "json",
						AddSource: true,
					},
					GRPCConfig: service.GRPCConfig{
						Server: service.GRPCServerConfig{
							Enable: true,
						},
					},
					HTTPConfig: service.HTTPConfig{
						EnableSwagger: true,
					},
					PICConfig: service.PICConfig{
						Serial: service.SerialConfig{
							Port:        "/dev/tty.usbserial",
							BaudRate:    9600,
							DataBits:    8,
							StopBits:    1.0,
							Parity:      "none",
							ReadTimeout: time.Second * 1,
						},
					},
				},
				expectedError: false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				cfgManager := mocks.NewFakeManager(t)
				tc.mock(cfgManager)

				systemService := NewSystemService(cfgManager)
				result, err := systemService.GetSystemConfig(ctx)

				if tc.expectedError {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
				cfgManager.AssertExpectations(t)
			})
		}
	})

	t.Run("test UpdateSystemConfig", func(t *testing.T) {
		tests := []struct {
			name          string
			params        service.UpdateSystemConfigParams
			mock          func(_ *mocks.FakeManager)
			expected      service.UpdateSystemConfigOutput
			expectedError bool
			errorType     error
		}{
			{
				name: "successful update config",
				params: service.UpdateSystemConfigParams{
					LogConfig: service.LogConfig{
						Level:     "debug",
						Format:    "json",
						AddSource: true,
					},
					GRPCConfig: service.GRPCConfig{
						Server: service.GRPCServerConfig{
							Enable: true,
						},
						Cloud: service.CloudConfig{
							Address: "localhost:50051",
						},
					},
					HTTPConfig: service.HTTPConfig{
						EnableSwagger: true,
					},
					PICConfig: service.PICConfig{
						Serial: service.SerialConfig{
							Port:        "/dev/ttyUSB0",
							BaudRate:    115200,
							DataBits:    8,
							StopBits:    1.0,
							Parity:      "none",
							ReadTimeout: time.Second * 1,
						},
					},
				},
				mock: func(cfgManager *mocks.FakeManager) {
					cfg := config.Config{
						Log: log.Config{
							Level:     "info",
							Format:    "json",
							AddSource: false,
						},
						GRPC: config.GRPCConfig{
							Server: config.GRPCServerConfig{
								Enable: true,
							},
						},
						HTTP: http.Config{
							EnableSwagger: false,
						},
						PIC: picserial.Config{
							Serial: serial.Config{
								Port:        "/dev/ttyUSB1",
								BaudRate:    9600,
								DataBits:    8,
								StopBits:    1.0,
								Parity:      "none",
								ReadTimeout: time.Second * 2,
							},
						},
					}
					cfgManager.EXPECT().GetConfig().Return(cfg)

					// Create a matcher function for SetConfig
					cfgManager.EXPECT().SaveConfig(mock.Anything, mock.Anything).Run(func(_ context.Context, cfg config.Config) {
						assert.Equal(t, "debug", cfg.Log.Level)
						assert.Equal(t, "json", cfg.Log.Format)
						assert.True(t, cfg.Log.AddSource)
						assert.True(t, cfg.GRPC.Server.Enable)
						assert.True(t, cfg.HTTP.EnableSwagger)
						assert.Equal(t, "/dev/ttyUSB0", cfg.PIC.Serial.Port)
						assert.Equal(t, 115200, cfg.PIC.Serial.BaudRate)
						assert.Equal(t, 8, cfg.PIC.Serial.DataBits)
						assert.Equal(t, 1.0, cfg.PIC.Serial.StopBits)
						assert.Equal(t, "none", cfg.PIC.Serial.Parity)
						assert.Equal(t, time.Second*1, cfg.PIC.Serial.ReadTimeout)
					}).Return(nil)
				},
				expected: service.UpdateSystemConfigOutput{
					LogConfig: service.LogConfig{
						Level:     "debug",
						Format:    "json",
						AddSource: true,
					},
					GRPCConfig: service.GRPCConfig{
						Server: service.GRPCServerConfig{
							Enable: true,
						},
						Cloud: service.CloudConfig{
							Address: "localhost:50051",
						},
					},
					HTTPConfig: service.HTTPConfig{
						EnableSwagger: true,
					},
					PICConfig: service.PICConfig{
						Serial: service.SerialConfig{
							Port:        "/dev/ttyUSB0",
							BaudRate:    115200,
							DataBits:    8,
							StopBits:    1.0,
							Parity:      "none",
							ReadTimeout: time.Second * 1,
						},
					},
				},
				expectedError: false,
			},
			{
				name: "validation error",
				params: service.UpdateSystemConfigParams{
					LogConfig: service.LogConfig{
						Level:     "invalid",
						Format:    "json",
						AddSource: true,
					},
					GRPCConfig: service.GRPCConfig{
						Cloud: service.CloudConfig{
							Address: "localhost:50051",
						},
					},
				},
				mock: func(cfgManager *mocks.FakeManager) {
					cfg := config.Config{
						Log: log.Config{
							Level:     "info",
							Format:    "json",
							AddSource: false,
						},
					}
					cfgManager.EXPECT().GetConfig().Return(cfg)

					// Mock SetConfig to return validation error
					cfgManager.EXPECT().SaveConfig(mock.Anything, mock.Anything).Return(config.ErrInvalidConfig)
				},
				expectedError: true,
				errorType:     ErrInvalidConfig,
			},
			{
				name: "set config error",
				params: service.UpdateSystemConfigParams{
					LogConfig: service.LogConfig{
						Level:  "info",
						Format: "json",
					},
				},
				mock: func(cfgManager *mocks.FakeManager) {
					cfg := config.Config{
						Log: log.Config{
							Level:  "info",
							Format: "json",
						},
					}
					cfgManager.EXPECT().GetConfig().Return(cfg)

					// Mock generic error during SetConfig
					cfgManager.EXPECT().SaveConfig(mock.Anything, mock.Anything).Return(assert.AnError)
				},
				expectedError: true,
			},
			{
				name: "save config error",
				params: service.UpdateSystemConfigParams{
					LogConfig: service.LogConfig{
						Level:  "info",
						Format: "json",
					},
				},
				mock: func(cfgManager *mocks.FakeManager) {
					cfg := config.Config{
						Log: log.Config{
							Level:  "info",
							Format: "json",
						},
					}
					cfgManager.EXPECT().GetConfig().Return(cfg)

					// SetConfig succeeds but SaveConfig fails
					cfgManager.EXPECT().SaveConfig(mock.Anything, mock.Anything).Return(assert.AnError)
				},
				expectedError: true,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				cfgManager := mocks.NewFakeManager(t)
				tc.mock(cfgManager)

				systemService := NewSystemService(cfgManager)
				result, err := systemService.UpdateSystemConfig(ctx, tc.params)

				if tc.expectedError {
					assert.Error(t, err)
					if tc.errorType != nil {
						assert.ErrorIs(t, err, tc.errorType)
					}
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
				cfgManager.AssertExpectations(t)
			})
		}
	})
}
