package apierr

import (
	"errors"
	"net/http"

	govalidator "github.com/go-playground/validator/v10"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/pkg/validator"
	"github.com/tbe-team/raybot/pkg/xerror"
)

// ErrorResponse is the error response for the API.
type ErrorResponse struct {
	gen.ErrorResponse

	// StatusCode is the status code for the error response.
	StatusCode int `json:"-"`
}

func New(err error) ErrorResponse {
	return errorToErrorResponse(err)
}

func errorToErrorResponse(err error) ErrorResponse {
	var xErr xerror.XError
	if errors.As(err, &xErr) {
		return ErrorResponse{
			ErrorResponse: gen.ErrorResponse{
				Code:    xErr.MsgID(),
				Message: xErr.Msg(),
			},
			StatusCode: xErr.Status().HTTPStatus(),
		}
	}

	var validationErrs govalidator.ValidationErrors
	if errors.As(err, &validationErrs) {
		details := make([]gen.FieldError, len(validationErrs))
		for i, fe := range validationErrs {
			details[i] = gen.FieldError{
				Field:   fe.Field(),
				Message: validator.ValidationErrorMessage(fe),
			}
		}

		return ErrorResponse{
			ErrorResponse: gen.ErrorResponse{
				Code:    "validationError",
				Message: "validation error",
				Details: &details,
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	return ErrorResponse{
		ErrorResponse: gen.ErrorResponse{
			Code:    "internalServerError",
			Message: "an unknown error occurred",
		},
		StatusCode: http.StatusInternalServerError,
	}
}
