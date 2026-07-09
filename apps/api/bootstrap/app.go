// apps/api/bootstrap/app.go
package bootstrap

import (
	"fmt"
	"net"

	"portfolio-ai/internal/auth/jwt"
	authgrpc "portfolio-ai/internal/auth/grpc"
	"portfolio-ai/internal/profile/repository"
	"portfolio-ai/internal/profile/service"
	profilegrpc "portfolio-ai/internal/profile/grpc"
	"portfolio-ai/pkg/config"
	"portfolio-ai/pkg/logger"
	pb "portfolio-ai/proto/auth"
	profilepb "portfolio-ai/proto/profile"

	"google.golang.org/grpc"
	gormdb "gorm.io/gorm"
)

type App struct {
	Config     *config.Config
	DB         *gormdb.DB
	GRPCServer *grpc.Server
}

// NewApp initializes configuration, logger, database, and gRPC server.
func NewApp() (*App, error) {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Initialize structured logger
	logger.Init(cfg.App.Env)
	logger.Info("Starting application", "app_name", cfg.App.Name, "env", cfg.App.Env)

	// 3. Connect to database
	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	logger.Info("Database connection established", "host", cfg.DB.Host, "database", cfg.DB.Name)

	// 4. Initialize JWT manager
	jwtManager := jwt.NewManager(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiry)
	logger.Info("JWT manager initialized", "expiry", cfg.Auth.JWTExpiry.String())

	// 5. Create gRPC server (with auth interceptor)
	grpcServer := NewGRPCServer(jwtManager)
	logger.Info("gRPC server initialized")

	// 6. Register Auth service (login without DB)
	authHandler := authgrpc.NewHandler(jwtManager, cfg.Auth.AdminUsername, cfg.Auth.AdminPassword, cfg.Auth.JWTExpiry)
	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	logger.Info("Auth service registered")

	// 7. Initialize Profile module
	profileRepo := repository.NewPostgresRepository(db)
	profileSvc := service.NewService(profileRepo)
	profileHnd := profilegrpc.NewHandler(profileSvc)
	profilepb.RegisterProfileServiceServer(grpcServer, profileHnd)
	logger.Info("Profile service registered")

	return &App{
		Config:     cfg,
		DB:         db,
		GRPCServer: grpcServer,
	}, nil
}

// Run starts the gRPC server on the configured port.
func (a *App) Run() error {
	addr := ":" + a.Config.App.Port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	logger.Info("gRPC server listening", "address", addr)
	return a.GRPCServer.Serve(lis)
}
