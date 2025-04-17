package http

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/apperrorcode"
)

type errorCodeHandler struct {
	apperrorcodeService apperrorcode.Service
}

func newErrorCodeHandler(apperrorcodeService apperrorcode.Service) *errorCodeHandler {
	return &errorCodeHandler{
		apperrorcodeService: apperrorcodeService,
	}
}

func (h errorCodeHandler) GetErrorCodes(ctx context.Context, _ gen.GetErrorCodesRequestObject) (gen.GetErrorCodesResponseObject, error) {
	errorCodes, err := h.apperrorcodeService.ListErrorCodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list error codes: %w", err)
	}

	items := make([]gen.ErrorCodeResponse, len(errorCodes))
	for i, errorCode := range errorCodes {
		items[i] = h.convertErrorCodeToResponse(errorCode)
	}

	return gen.GetErrorCodes200JSONResponse(items), nil
}

func (errorCodeHandler) convertErrorCodeToResponse(errorCode apperrorcode.ErrorCode) gen.ErrorCodeResponse {
	return gen.ErrorCodeResponse{
		Code:    errorCode.Code,
		Message: errorCode.Message,
	}
}
