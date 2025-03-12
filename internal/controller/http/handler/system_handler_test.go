package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/service/mocks"
)

func TestSystemHandler(t *testing.T) {
	t.Run("test get system config", func(t *testing.T) {
		tests := []struct {
			name           string
			mockSystemSvc  func(svc *mocks.FakeSystemService)
			expectedStatus int
			expectedBody   *gen.SystemConfigResponse
		}{
			{
				name: "successful get system config",
				mockSystemSvc: func(svc *mocks.FakeSystemService) {
					svc.EXPECT().GetSystemConfig(context.Background()).Return(
						service.GetSystemConfigOutput{
							LogConfig: service.LogConfig{
								Level:     "info",
								Format:    "json",
								AddSource: true,
							},
							GRPCConfig: service.GRPCConfig{
								Port: 50051,
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
									ReadTimeout: 5 * time.Second,
								},
							},
						}, nil)
				},
				expectedStatus: http.StatusOK,
				expectedBody: &gen.SystemConfigResponse{
					Log: gen.LogConfig{
						Level:     "info",
						Format:    "json",
						AddSource: true,
					},
					Grpc: gen.GRPCConfig{
						Port: 50051,
					},
					Http: gen.HTTPConfig{
						EnableSwagger: true,
					},
					Pic: gen.PicConfig{
						Serial: gen.SerialConfig{
							Port:        "/dev/ttyUSB0",
							BaudRate:    115200,
							DataBits:    8,
							StopBits:    1.0,
							Parity:      "none",
							ReadTimeout: 5.0,
						},
					},
				},
			},
			{
				name: "service error",
				mockSystemSvc: func(svc *mocks.FakeSystemService) {
					svc.EXPECT().GetSystemConfig(context.Background()).
						Return(service.GetSystemConfigOutput{}, errors.New("internal error"))
				},
				expectedStatus: http.StatusInternalServerError,
				expectedBody:   nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				systemSvc := mocks.NewFakeSystemService(t)
				tc.mockSystemSvc(systemSvc)

				handler := systemHandler{
					systemService: systemSvc,
				}

				// Create a new request
				req := httptest.NewRequest("GET", "/system/config", nil)
				rec := httptest.NewRecorder()

				result, err := handler.GetSystemConfig(req.Context(), gen.GetSystemConfigRequestObject{})

				// Process the result to simulate middleware behavior
				if err != nil {
					rec.WriteHeader(http.StatusInternalServerError)
					if _, err := rec.Write([]byte(err.Error())); err != nil {
						assert.Fail(t, "failed to write response body")
					}
				} else if response, ok := result.(gen.GetSystemConfig200JSONResponse); ok {
					data, _ := json.Marshal(response)
					rec.Header().Set("Content-Type", "application/json")
					rec.WriteHeader(http.StatusOK)
					if _, err := rec.Write(data); err != nil {
						assert.Fail(t, "failed to write response body")
					}
				}

				// Assertions
				assert.Equal(t, tc.expectedStatus, rec.Code)

				if tc.expectedBody != nil {
					var responseBody gen.SystemConfigResponse
					err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
					require.NoError(t, err)
					assert.Equal(t, *tc.expectedBody, responseBody)
				}

				systemSvc.AssertExpectations(t)
			})
		}
	})

	t.Run("test update system config", func(t *testing.T) {
		tests := []struct {
			name           string
			requestBody    gen.SystemConfigRequest
			mockSystemSvc  func(svc *mocks.FakeSystemService)
			expectedStatus int
			expectedBody   *gen.SystemConfigResponse
		}{
			{
				name: "successful update system config",
				requestBody: gen.SystemConfigRequest{
					Log: gen.LogConfig{
						Level:     "debug",
						Format:    "json",
						AddSource: true,
					},
					Grpc: gen.GRPCConfig{
						Port: 50051,
					},
					Http: gen.HTTPConfig{
						EnableSwagger: true,
					},
					Pic: gen.PicConfig{
						Serial: gen.SerialConfig{
							Port:        "/dev/ttyUSB0",
							BaudRate:    115200,
							DataBits:    8,
							StopBits:    1.0,
							Parity:      "none",
							ReadTimeout: 5.0,
						},
					},
				},
				mockSystemSvc: func(svc *mocks.FakeSystemService) {
					svc.EXPECT().UpdateSystemConfig(
						context.Background(),
						mock.MatchedBy(func(params service.UpdateSystemConfigParams) bool {
							return params.LogConfig.Level == "debug" &&
								params.LogConfig.Format == "json" &&
								params.LogConfig.AddSource == true &&
								params.GRPCConfig.Port == 50051 &&
								params.HTTPConfig.EnableSwagger == true &&
								params.PICConfig.Serial.Port == "/dev/ttyUSB0" &&
								params.PICConfig.Serial.BaudRate == 115200 &&
								params.PICConfig.Serial.DataBits == 8 &&
								params.PICConfig.Serial.StopBits == 1.0 &&
								params.PICConfig.Serial.Parity == "none" &&
								params.PICConfig.Serial.ReadTimeout == 5*time.Second
						}),
					).Return(
						service.UpdateSystemConfigOutput{
							LogConfig: service.LogConfig{
								Level:     "debug",
								Format:    "json",
								AddSource: true,
							},
							GRPCConfig: service.GRPCConfig{
								Port: 50051,
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
									ReadTimeout: 5 * time.Second,
								},
							},
						}, nil)
				},
				expectedStatus: http.StatusOK,
				expectedBody: &gen.SystemConfigResponse{
					Log: gen.LogConfig{
						Level:     "debug",
						Format:    "json",
						AddSource: true,
					},
					Grpc: gen.GRPCConfig{
						Port: 50051,
					},
					Http: gen.HTTPConfig{
						EnableSwagger: true,
					},
					Pic: gen.PicConfig{
						Serial: gen.SerialConfig{
							Port:        "/dev/ttyUSB0",
							BaudRate:    115200,
							DataBits:    8,
							StopBits:    1.0,
							Parity:      "none",
							ReadTimeout: 5.0,
						},
					},
				},
			},
			{
				name: "service error",
				requestBody: gen.SystemConfigRequest{
					Log: gen.LogConfig{
						Level:     "invalid",
						Format:    "json",
						AddSource: true,
					},
					Grpc: gen.GRPCConfig{
						Port: 50051,
					},
					Http: gen.HTTPConfig{
						EnableSwagger: true,
					},
					Pic: gen.PicConfig{
						Serial: gen.SerialConfig{
							Port:        "/dev/ttyUSB0",
							BaudRate:    115200,
							DataBits:    8,
							StopBits:    1.0,
							Parity:      "none",
							ReadTimeout: 5.0,
						},
					},
				},
				mockSystemSvc: func(svc *mocks.FakeSystemService) {
					svc.EXPECT().UpdateSystemConfig(context.Background(), mock.Anything).Return(
						service.UpdateSystemConfigOutput{}, errors.New("validation failed: invalid log level"))
				},
				expectedStatus: http.StatusInternalServerError,
				expectedBody:   nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				systemSvc := mocks.NewFakeSystemService(t)
				tc.mockSystemSvc(systemSvc)

				handler := systemHandler{
					systemService: systemSvc,
				}

				// Create a new request with body
				body, err := json.Marshal(tc.requestBody)
				require.NoError(t, err)
				req := httptest.NewRequest("PUT", "/system/config", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				// Create the request object directly
				requestObj := gen.UpdateSystemConfigRequestObject{
					Body: &tc.requestBody,
				}

				// Call the method directly
				result, err := handler.UpdateSystemConfig(context.Background(), requestObj)

				// Process the result as middleware would
				if err != nil {
					rec.WriteHeader(http.StatusInternalServerError)
					if _, err := rec.Write([]byte(err.Error())); err != nil {
						assert.Fail(t, "failed to write response body")
					}
				} else if response, ok := result.(gen.UpdateSystemConfig200JSONResponse); ok {
					data, _ := json.Marshal(response)
					rec.Header().Set("Content-Type", "application/json")
					rec.WriteHeader(http.StatusOK)
					if _, err := rec.Write(data); err != nil {
						assert.Fail(t, "failed to write response body")
					}
				}

				// Assertions
				assert.Equal(t, tc.expectedStatus, rec.Code)

				if tc.expectedBody != nil {
					var responseBody gen.SystemConfigResponse
					err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
					require.NoError(t, err)
					assert.Equal(t, *tc.expectedBody, responseBody)
				}

				systemSvc.AssertExpectations(t)
			})
		}
	})

	t.Run("test restart application", func(t *testing.T) {
		tests := []struct {
			name           string
			mockSystemSvc  func(_ *mocks.FakeSystemService)
			expectedStatus int
		}{
			{
				name: "successful restart",
				mockSystemSvc: func(svc *mocks.FakeSystemService) {
					svc.EXPECT().
						RestartApplication(context.Background()).
						Return(nil)
				},
				expectedStatus: http.StatusNoContent,
			},
			{
				name: "restart error",
				mockSystemSvc: func(svc *mocks.FakeSystemService) {
					svc.EXPECT().
						RestartApplication(context.Background()).
						Return(assert.AnError)
				},
				expectedStatus: http.StatusInternalServerError,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				systemSvc := mocks.NewFakeSystemService(t)
				tc.mockSystemSvc(systemSvc)

				handler := systemHandler{
					systemService: systemSvc,
				}

				// Create a new request
				req := httptest.NewRequest("POST", "/system/restart", nil)
				rec := httptest.NewRecorder()

				// Call the method directly
				result, err := handler.RestartApplication(req.Context(), gen.RestartApplicationRequestObject{})

				// Process the result to simulate middleware behavior
				if err != nil {
					rec.WriteHeader(http.StatusInternalServerError)
					if _, err := rec.Write([]byte(err.Error())); err != nil {
						assert.Fail(t, "failed to write response body")
					}
				} else if _, ok := result.(gen.RestartApplication204Response); ok {
					rec.WriteHeader(http.StatusNoContent)
				}

				// Assertions
				assert.Equal(t, tc.expectedStatus, rec.Code)
				systemSvc.AssertExpectations(t)
			})
		}
	})
}
