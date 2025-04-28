package grpcerr

import (
	"errors"
	"testing"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tbe-team/raybot/pkg/xerror"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		inputErr       error
		expectedCode   codes.Code
		expectedErrMsg string
	}{
		{
			name:           "XError with unauthorized status",
			inputErr:       xerror.Unauthorized(nil, "auth.unauthorized", "unauthorized"),
			expectedCode:   codes.Unauthenticated,
			expectedErrMsg: "unauthorized",
		},
		{
			name:           "XError with not found status",
			inputErr:       xerror.NotFound(nil, "entity.notFound", "not found"),
			expectedCode:   codes.NotFound,
			expectedErrMsg: "not found",
		},
		{
			name:           "Validation error",
			inputErr:       govalidator.ValidationErrors{},
			expectedCode:   codes.InvalidArgument,
			expectedErrMsg: "validation error",
		},
		{
			name:           "Generic error",
			inputErr:       errors.New("some error"),
			expectedCode:   codes.Internal,
			expectedErrMsg: "internal server error: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.inputErr)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, st.Code())
			assert.Contains(t, err.Error(), tt.expectedErrMsg)
		})
	}
}

func TestXErrorStatusToGRPCCode(t *testing.T) {
	tests := []struct {
		name         string
		status       xerror.Status
		expectedCode codes.Code
	}{
		{"Unauthorized", xerror.StatusUnauthorized, codes.Unauthenticated},
		{"Forbidden", xerror.StatusForbidden, codes.PermissionDenied},
		{"NotFound", xerror.StatusNotFound, codes.NotFound},
		{"UnprocessableEntity", xerror.StatusUnprocessableEntity, codes.FailedPrecondition},
		{"Conflict", xerror.StatusConflict, codes.AlreadyExists},
		{"TooManyRequests", xerror.StatusTooManyRequests, codes.ResourceExhausted},
		{"BadRequest", xerror.StatusBadRequest, codes.InvalidArgument},
		{"ValidationFailed", xerror.StatusValidationFailed, codes.InvalidArgument},
		{"Unknown", xerror.StatusUnknown, codes.Internal},
		{"InternalServerError", xerror.StatusInternalServerError, codes.Internal},
		{"Timeout", xerror.StatusTimeout, codes.DeadlineExceeded},
		{"NotImplemented", xerror.StatusNotImplemented, codes.Unimplemented},
		{"BadGateway", xerror.StatusBadGateway, codes.Unavailable},
		{"Default", xerror.StatusInternalServerError, codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := XErrorStatusToGRPCCode(tt.status)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}
