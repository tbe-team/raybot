package interceptor

import (
	"context"

	"google.golang.org/grpc"

	"github.com/tbe-team/raybot/internal/handlers/cloud/grpcerr"
)

func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, grpcerr.New(err)
		}
		return resp, nil
	}
}

func StreamErrorInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		_ *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, stream)
		if err != nil {
			return grpcerr.New(err)
		}
		return nil
	}
}
