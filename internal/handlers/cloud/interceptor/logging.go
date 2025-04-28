package interceptor

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

func UnaryLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(interceptorLogger(log))
}

func StreamLoggingInterceptor(log *slog.Logger) grpc.StreamServerInterceptor {
	return logging.StreamServerInterceptor(interceptorLogger(log))
}

func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
