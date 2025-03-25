package picserial

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalMessageType(t *testing.T) {
	tests := []struct {
		name      string
		msg       []byte
		shouldErr bool
		want      messageType
	}{
		{
			name:      "sync state",
			msg:       []byte(`{"type": 0}`),
			shouldErr: false,
			want:      messageTypeSyncState,
		},
		{
			name:      "ack message type",
			msg:       []byte(`{"type": 1}`),
			shouldErr: false,
			want:      messageTypeACK,
		},
		{
			name:      "invalid message type",
			msg:       []byte(`{"type": 2}`),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var temp struct {
				Type messageType `json:"type"`
			}
			err := json.Unmarshal(tt.msg, &temp)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("failed to unmarshal message type: %v", err)
				}
				if temp.Type != tt.want {
					t.Errorf("got %v, want %v", temp.Type, tt.want)
				}
			}
		})
	}
}
