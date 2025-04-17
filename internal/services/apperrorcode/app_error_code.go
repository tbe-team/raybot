package apperrorcode

import "context"

type Service interface {
	ListErrorCodes(ctx context.Context) ([]ErrorCode, error)
}
