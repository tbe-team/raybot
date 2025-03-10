package handler_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/controller/picserial/handler"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/service/mocks"
	"github.com/tbe-team/raybot/pkg/log"
)

func TestACKStatusUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     handler.ACKStatus
		wantErr  bool
	}{
		{
			name:     "success status",
			jsonData: `1`,
			want:     handler.ACKStatusSuccess,
			wantErr:  false,
		},
		{
			name:     "failure status",
			jsonData: `0`,
			want:     handler.ACKStatusFailure,
			wantErr:  false,
		},
		{
			name:     "invalid status",
			jsonData: `2`,
			wantErr:  true,
		},
		{
			name:     "non-numeric value",
			jsonData: `"abc"`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status handler.ACKStatus
			err := json.Unmarshal([]byte(tt.jsonData), &status)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, status)
			}
		})
	}
}

func TestCommandACKHandler_Handle(t *testing.T) {
	log := log.NewNopLogger()
	ctx := context.Background()

	tests := []struct {
		name          string
		msg           handler.CommandACKMessage
		mockReturn    error
		expectedError bool
	}{
		{
			name: "process success",
			msg: handler.CommandACKMessage{
				ID:     "1",
				Status: handler.ACKStatusSuccess,
			},
			mockReturn:    nil,
			expectedError: false,
		},
		{
			name: "process failure",
			msg: handler.CommandACKMessage{
				ID:     "2",
				Status: handler.ACKStatusSuccess,
			},
			mockReturn:    assert.AnError,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			picService := mocks.NewFakePICService(t)
			commandACKHandler := handler.NewCommandACKHandler(picService, log)

			params := service.ProcessSerialCommandACKParams{
				ID:      tc.msg.ID,
				Success: tc.msg.Status == handler.ACKStatusSuccess,
			}
			picService.EXPECT().ProcessSerialCommandACK(ctx, params).Return(tc.mockReturn)

			commandACKHandler.Handle(ctx, tc.msg)

			if tc.expectedError {
				assert.Error(t, tc.mockReturn)
			} else {
				assert.NoError(t, tc.mockReturn)
			}
		})
	}
}
