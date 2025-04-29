package cloudtest

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryTunneledAuthnInterceptor(headerKey, requireToken string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		token, err := getTokenFromContext(ctx, headerKey)
		if err != nil {
			return nil, fmt.Errorf("get token from context: %w", err)
		}

		if token != requireToken {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}

func StreamTunneledAuthnInterceptor(headerKey, requireToken string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		_ *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		token, err := getTokenFromContext(stream.Context(), headerKey)
		if err != nil {
			return fmt.Errorf("get token from context: %w", err)
		}

		if token != requireToken {
			return status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return handler(srv, stream)
	}
}

func getTokenFromContext(ctx context.Context, headerKey string) (string, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	tokens := metadata.Get(headerKey)
	if len(tokens) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing credentials")
	}

	return tokens[0], nil
}
