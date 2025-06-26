package cloud

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const RequestIDKey = "request-id"

// GetRequestIDFromContext retrieves the request ID from the context metadata.
// If the request ID is not present, it returns nil.
func GetRequestIDFromContext(ctx context.Context) *string {
	var requestID *string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if reqID := md.Get(RequestIDKey); len(reqID) > 0 {
			requestID = &reqID[0]
		}
	}

	return requestID
}
