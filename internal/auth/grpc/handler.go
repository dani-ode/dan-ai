// internal/auth/grpc/handler.go
package grpc

import (
	"context"
	"time"

	"dan-ai/internal/auth/jwt"
	pb "dan-ai/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler implements the AuthService gRPC server.
// It validates credentials against environment variables (no database).
type Handler struct {
	pb.UnimplementedAuthServiceServer
	jwtManager *jwt.Manager
	username   string
	password   string
	expiry     time.Duration
}

// NewHandler creates a new auth gRPC handler.
func NewHandler(jwtManager *jwt.Manager, username, password string, expiry time.Duration) *Handler {
	return &Handler{
		jwtManager: jwtManager,
		username:   username,
		password:   password,
		expiry:     expiry,
	}
}

// Login validates credentials from env and returns a JWT token.
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	if req.GetUsername() != h.username || req.GetPassword() != h.password {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, err := h.jwtManager.GenerateToken(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	expiresAt := time.Now().Add(h.expiry).Format(time.RFC3339)

	return &pb.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
