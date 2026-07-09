// apps/api/bootstrap/grpc.go
package bootstrap

import (
	"portfolio-ai/internal/auth/jwt"
	pkggrpc "portfolio-ai/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// NewGRPCServer creates a gRPC server with interceptors and registers the health service.
func NewGRPCServer(jwtManager *jwt.Manager) *grpc.Server {
	srv := pkggrpc.NewServer(jwtManager)

	// Register gRPC health check service
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(srv, healthServer)

	return srv
}
