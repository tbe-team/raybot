package grpcerr

import (
	"errors"
	"fmt"

	govalidator "github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tbe-team/raybot/pkg/xerror"
)

func New(err error) error {
	var xErr xerror.XError
	if errors.As(err, &xErr) {
		return status.Error(XErrorStatusToGRPCCode(xErr.Status()), xErr.Msg())
	}

	var validationErrs govalidator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return status.Error(codes.InvalidArgument, "validation error")
	}

	grpcErr := status.Error(codes.Internal, "internal server error")
	return fmt.Errorf("%w: %s", grpcErr, err)
}

func XErrorStatusToGRPCCode(status xerror.Status) codes.Code {
	switch status {
	case xerror.StatusUnauthorized:
		return codes.Unauthenticated
	case xerror.StatusForbidden:
		return codes.PermissionDenied
	case xerror.StatusNotFound:
		return codes.NotFound
	case xerror.StatusUnprocessableEntity:
		return codes.FailedPrecondition
	case xerror.StatusConflict:
		return codes.AlreadyExists
	case xerror.StatusTooManyRequests:
		return codes.ResourceExhausted
	case xerror.StatusBadRequest:
		return codes.InvalidArgument
	case xerror.StatusValidationFailed:
		return codes.InvalidArgument
	case xerror.StatusUnknown, xerror.StatusInternalServerError:
		return codes.Internal
	case xerror.StatusTimeout:
		return codes.DeadlineExceeded
	case xerror.StatusNotImplemented:
		return codes.Unimplemented
	case xerror.StatusBadGateway:
		return codes.Unavailable
	default:
		return codes.Internal
	}
}
