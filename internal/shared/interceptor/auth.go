// internal/shared/interceptor/auth.go
package interceptor

import (
	"context"
	"strings"

	"portfolio-ai/internal/auth/jwt"
	"portfolio-ai/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// publicMethods are gRPC methods that do NOT require authentication.
var publicMethods = map[string]bool{
	"/auth.AuthService/Login":         true,
	"/profile.ProfileService/GetProfile": true,
	"/grpc.health.v1.Health/Check":    true,
	"/grpc.health.v1.Health/Watch":    true,
	"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":   true,
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
}

// UnaryAuth returns a unary server interceptor that validates JWT tokens.
// Public methods (Login, Health, Reflection) are exempt from auth.
func UnaryAuth(jwtManager *jwt.Manager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		claims, err := extractAndValidate(ctx, jwtManager)
		if err != nil {
			logger.Warn("Auth failed", "method", info.FullMethod, "error", err)
			return nil, err
		}

		// Store claims in context for downstream handlers.
		ctx = context.WithValue(ctx, claimsKey{}, claims)
		return handler(ctx, req)
	}
}

// StreamAuth returns a stream server interceptor that validates JWT tokens.
func StreamAuth(jwtManager *jwt.Manager) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if publicMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		_, err := extractAndValidate(ss.Context(), jwtManager)
		if err != nil {
			logger.Warn("Auth failed", "method", info.FullMethod, "error", err)
			return err
		}

		return handler(srv, ss)
	}
}

// extractAndValidate extracts the bearer token from gRPC metadata and validates it.
func extractAndValidate(ctx context.Context, jwtManager *jwt.Manager) (*jwt.Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	token := values[0]
	if !strings.HasPrefix(token, "Bearer ") {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization format, use: Bearer <token>")
	}

	tokenStr := strings.TrimPrefix(token, "Bearer ")
	claims, err := jwtManager.ValidateToken(tokenStr)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return claims, nil
}

// claimsKey is the context key for JWT claims.
type claimsKey struct{}

// ClaimsFromContext extracts JWT claims from the context.
func ClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*jwt.Claims)
	return claims, ok
}
