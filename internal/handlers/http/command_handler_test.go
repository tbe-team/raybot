package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/command"
	commandmocks "github.com/tbe-team/raybot/internal/services/command/mocks"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/ptr"
)

func TestCommandHandler_GetCommandById(t *testing.T) {
	t.Run("Should get successfully", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().GetCommandByID(mock.Anything,
			mock.MatchedBy(
				func(params command.GetCommandByIDParams) bool {
					return params.CommandID == 123
				},
			),
		).Return(validCommand, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commands/123", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Should not able to get command if fetching failed", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().GetCommandByID(mock.Anything,
			mock.MatchedBy(
				func(params command.GetCommandByIDParams) bool {
					return params.CommandID == 123
				},
			),
		).Return(command.Command{}, errors.New("fetching failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commands/123", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestCommandHandler_GetCurrentProcessingCommand(t *testing.T) {
	t.Run("Should get successfully", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().GetCurrentProcessingCommand(mock.Anything).Return(validCommand, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commands/processing", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Should not able to get current processing command if fetching failed", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().GetCurrentProcessingCommand(mock.Anything).
			Return(command.Command{}, errors.New("fetching failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commands/processing", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestCommandHandler_ListCommands(t *testing.T) {
	t.Run("Should get successfully", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().ListCommands(mock.Anything,
			mock.MatchedBy(
				func(params command.ListCommandsParams) bool {
					return params.PagingParams.Page == 3 &&
						params.PagingParams.PageSize == 32 &&
						len(params.Statuses) == 1 &&
						params.Statuses[0] == command.StatusSucceeded &&
						params.Sorts[0].Col == "created_at" &&
						params.Sorts[0].Order == "DESC"
				},
			),
		).Return(paging.List[command.Command]{
			Items:      []command.Command{validCommand},
			TotalItems: 1,
		}, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commands?page=3&pageSize=32&statuses=SUCCEEDED&sorts=-created_at", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		res := MustDecodeJSON[gen.ListCommands200JSONResponse](t, rec.Body)
		require.Equal(t, 1, res.TotalItems)
		require.Equal(t, validCommand.ID, int64(res.Items[0].Id))
	})

	t.Run("Should use default paging params if not provided", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().ListCommands(mock.Anything,
			mock.MatchedBy(
				func(params command.ListCommandsParams) bool {
					return params.PagingParams.Page == 1 && params.PagingParams.PageSize == 10
				},
			),
		).Return(paging.List[command.Command]{
			Items:      []command.Command{validCommand},
			TotalItems: 1,
		}, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commands", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Should not able to list commands if fetching failed", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().ListCommands(mock.Anything, mock.Anything).
			Return(paging.List[command.Command]{}, errors.New("fetching failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commands", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestCommandHandler_CreateCommand(t *testing.T) {
	t.Run("Should create command successfully", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().CreateCommand(mock.Anything,
			mock.MatchedBy(
				func(params command.CreateCommandParams) bool {
					i, ok := params.Inputs.(*command.CargoLiftInputs)
					return params.Source == command.SourceApp &&
						ok &&
						i.Position == 12 &&
						i.MotorSpeed == 12
				},
			),
		).Return(validCommand, nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		i := gen.CommandInputs{}
		err := i.FromCargoLiftInputs(gen.CargoLiftInputs{
			Position:   12,
			MotorSpeed: 12,
		})
		require.NoError(t, err)

		body := gen.CreateCommandRequest{
			Type:   "CARGO_LIFT",
			Inputs: i,
		}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/commands", bytes.NewBuffer(jsonBody))
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("Should not able to create command if validation fail", func(t *testing.T) {
		t.Run("Type is invalid", func(t *testing.T) {
			h := SetupAPITestHandler(t)

			body := gen.CreateCommandRequest{
				Type: "INVALID_TYPE",
			}
			jsonBody, err := json.Marshal(body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/commands", bytes.NewBuffer(jsonBody))
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("Inputs is invalid", func(t *testing.T) {
			h := SetupAPITestHandler(t)

			body := gen.CreateCommandRequest{
				Type:   "CARGO_LIFT",
				Inputs: gen.CommandInputs{},
			}
			jsonBody, err := json.Marshal(body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/commands", bytes.NewBuffer(jsonBody))
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("Should not able to create command if creating failed", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().CreateCommand(mock.Anything, mock.Anything).
			Return(command.Command{}, errors.New("creating failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		i := gen.CommandInputs{}
		err := i.FromCargoLiftInputs(gen.CargoLiftInputs{
			Position:   12,
			MotorSpeed: 12,
		})
		require.NoError(t, err)

		body := gen.CreateCommandRequest{
			Type:   "CARGO_LIFT",
			Inputs: i,
		}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/commands", bytes.NewBuffer(jsonBody))
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestCommandHandler_DeleteCommandById(t *testing.T) {
	t.Run("Should delete command successfully", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().DeleteCommandByID(mock.Anything,
			mock.MatchedBy(
				func(params command.DeleteCommandByIDParams) bool {
					return params.CommandID == 123
				},
			),
		).Return(nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/commands/123", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Should not able to delete command if deleting failed", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().DeleteCommandByID(mock.Anything, mock.Anything).
			Return(errors.New("deleting failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/commands/123", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestCommandHandler_CancelCurrentProcessingCommand(t *testing.T) {
	t.Run("Should cancel current processing command successfully", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().CancelCurrentProcessingCommand(mock.Anything).Return(nil)

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/commands/processing/cancel", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Should not able to cancel current processing command if fetching failed", func(t *testing.T) {
		commandService := commandmocks.NewFakeService(t)
		commandService.EXPECT().CancelCurrentProcessingCommand(mock.Anything).
			Return(errors.New("fetching failed"))

		h := SetupAPITestHandler(t, func(hs *Service) {
			hs.commandService = commandService
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/commands/processing/cancel", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

var validCommand = command.Command{
	ID:          1,
	Type:        command.CommandTypeStopMovement,
	Status:      command.StatusSucceeded,
	Source:      command.SourceApp,
	Inputs:      &command.StopMovementInputs{},
	Outputs:     &command.StopMovementOutputs{},
	StartedAt:   ptr.New(time.Now()),
	CompletedAt: ptr.New(time.Now()),
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}
