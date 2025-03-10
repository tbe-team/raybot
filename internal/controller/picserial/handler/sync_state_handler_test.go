package handler_test

import (
	"encoding/json"
	"testing"

	"github.com/tbe-team/raybot/internal/controller/picserial/handler"
)

func TestUnmarshalSyncStateType(t *testing.T) {
	tests := []struct {
		name      string
		msg       []byte
		shouldErr bool
		want      handler.SyncStateType
	}{
		{
			name:      "sync state type battery",
			msg:       []byte(`{"state_type": 0}`),
			shouldErr: false,
			want:      handler.SyncStateTypeBattery,
		},
		{
			name:      "sync state type charge",
			msg:       []byte(`{"state_type": 1}`),
			shouldErr: false,
			want:      handler.SyncStateTypeCharge,
		},
		{
			name:      "sync state type discharge",
			msg:       []byte(`{"state_type": 2}`),
			shouldErr: false,
			want:      handler.SyncStateTypeDischarge,
		},
		{
			name:      "sync state type distance sensor",
			msg:       []byte(`{"state_type": 3}`),
			shouldErr: false,
			want:      handler.SyncStateTypeDistanceSensor,
		},
		{
			name:      "sync state type lift motor",
			msg:       []byte(`{"state_type": 4}`),
			shouldErr: false,
			want:      handler.SyncStateTypeLiftMotor,
		},
		{
			name:      "sync state type drive motor",
			msg:       []byte(`{"state_type": 5}`),
			shouldErr: false,
			want:      handler.SyncStateTypeDriveMotor,
		},
		{
			name:      "invalid sync state type",
			msg:       []byte(`{"state_type": 6}`),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var temp struct {
				StateType handler.SyncStateType `json:"state_type"`
			}
			err := json.Unmarshal(tt.msg, &temp)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if temp.StateType != tt.want {
					t.Errorf("expected %v, got %v", tt.want, temp.StateType)
				}
			}
		})
	}
}
