package interceptor

import (
	"log/slog"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryRecoveryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(
		recovery.WithRecoveryHandler(recoveryHandler(log)),
	)
}

func StreamRecoveryInterceptor(log *slog.Logger) grpc.StreamServerInterceptor {
	return recovery.StreamServerInterceptor(
		recovery.WithRecoveryHandler(recoveryHandler(log)),
	)
}

func recoveryHandler(log *slog.Logger) func(p any) error {
	return func(p any) error {
		log.Error("recovered from panic",
			slog.Any("panic", p),
			slog.String("stack", string(debug.Stack())),
		)

		return status.Errorf(codes.Internal, "internal server error")
	}
}
