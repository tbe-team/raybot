package validator

import (
	"fmt"
	"testing"
)

type testEnum uint8

func (s testEnum) Validate() error {
	switch s {
	case testEnumValue1:
		return nil
	default:
		return fmt.Errorf("invalid enum value: %d", s)
	}
}

const (
	testEnumValue1 testEnum = iota
)

func TestValidateEnum(t *testing.T) {
	validator := New()
	type schema struct {
		Enum testEnum `validate:"enum"`
	}
	tests := []struct {
		name    string
		schema  schema
		wantErr bool
	}{
		{
			name:    "valid enum",
			schema:  schema{Enum: 0},
			wantErr: false,
		},
		{
			name:    "invalid enum",
			schema:  schema{Enum: 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.schema)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
