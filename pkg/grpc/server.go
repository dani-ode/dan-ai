// pkg/grpc/server.go
package grpc

import (
	"dan-ai/internal/auth/jwt"
	"dan-ai/internal/shared/interceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewServer creates a gRPC server with logging, recovery, and auth interceptors.
func NewServer(jwtManager *jwt.Manager) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryRecovery(),
			interceptor.UnaryLogger(),
			interceptor.UnaryAuth(jwtManager),
		),
		grpc.ChainStreamInterceptor(
			interceptor.StreamRecovery(),
			interceptor.StreamLogger(),
			interceptor.StreamAuth(jwtManager),
		),
	)

	// Enable server reflection for tools like grpcurl
	reflection.Register(srv)

	return srv
}
