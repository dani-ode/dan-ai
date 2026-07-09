// internal/shared/interceptor/logger.go
package interceptor

import (
	"context"
	"time"

	"portfolio-ai/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryLogger returns a unary server interceptor that logs each gRPC call.
func UnaryLogger() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		st, _ := status.FromError(err)
		logger.Info("gRPC Request",
			"method", info.FullMethod,
			"code", st.Code().String(),
			"latency", time.Since(start).String(),
		)

		return resp, err
	}
}

// StreamLogger returns a stream server interceptor that logs each gRPC stream.
func StreamLogger() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		err := handler(srv, ss)

		st, _ := status.FromError(err)
		logger.Info("gRPC Stream",
			"method", info.FullMethod,
			"code", st.Code().String(),
			"latency", time.Since(start).String(),
		)

		return err
	}
}
