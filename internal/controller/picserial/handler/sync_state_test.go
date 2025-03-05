package handler

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalSyncStateType(t *testing.T) {
	tests := []struct {
		name      string
		msg       []byte
		shouldErr bool
		want      syncStateType
	}{
		{
			name:      "sync state type battery",
			msg:       []byte(`{"state_type": 0}`),
			shouldErr: false,
			want:      syncStateTypeBattery,
		},
		{
			name:      "sync state type charge",
			msg:       []byte(`{"state_type": 1}`),
			shouldErr: false,
			want:      syncStateTypeCharge,
		},
		{
			name:      "sync state type discharge",
			msg:       []byte(`{"state_type": 2}`),
			shouldErr: false,
			want:      syncStateTypeDischarge,
		},
		{
			name:      "sync state type distance sensor",
			msg:       []byte(`{"state_type": 3}`),
			shouldErr: false,
			want:      syncStateTypeDistanceSensor,
		},
		{
			name:      "sync state type lift motor",
			msg:       []byte(`{"state_type": 4}`),
			shouldErr: false,
			want:      syncStateTypeLiftMotor,
		},
		{
			name:      "sync state type drive motor",
			msg:       []byte(`{"state_type": 5}`),
			shouldErr: false,
			want:      syncStateTypeDriveMotor,
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
				StateType syncStateType `json:"state_type"`
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
